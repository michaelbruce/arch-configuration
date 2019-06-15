package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sysconf/boot"
	"sysconf/install"
	"sysconf/pacman"
)

func installOperation(target string) {
	fmt.Printf("installing Arch Linux on %v...\n", target)
	err := install.CheckCapacity(target)

	if err != nil {
		log.Fatal(err)
	}

	err = install.CheckUEFI()

	if err != nil {
		log.Fatal(err)
	}

	// setup disk table and partitions
	out, err := exec.Command("parted", "--script",
		target, "mklabel", "gpt",
		"mkpart", "ESP", "fat32", "1M", "512M",
		"set", "1", "boot", "on",
		"mkpart", "system", "ext4", "512M", "100%").Output()

	fmt.Println(out)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	out, err = exec.Command("mkfs.fat", "-F32", target+"1").Output()

	fmt.Println(out)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	out, err = exec.Command("mkfs.ext4", "-F", target+"2").Output()

	fmt.Println(out)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	out, err = exec.Command("mount", target+"2", "/mnt").Output()

	fmt.Println(out)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err = os.Mkdir("/mnt/boot", 0755); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	out, err = exec.Command("mount", target+"1", "/mnt/boot").Output()

	fmt.Println(out)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// install pacman, the package manager
	out, err = exec.Command("pacstrap", "-i", "/mnt", "base", "base-devel", "--noconfirm").Output()

	fmt.Println(out)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// create file system table
	out, err = exec.Command("genfstab", "-U", "/mnt", ">", "/mnt/etc/fstab").Output()

	fmt.Println(out)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("base installation complete, chroot into the system and run ./sysconf -update")
}

func update(packages pacman.Packages) {
	fmt.Println("updating system configuration...")

	if !pacman.Exists() {
		os.Exit(1)
	}

	fmt.Printf("packages to include %v\n", packages)

	installedPackages := pacman.InstalledPackages()

	fmt.Printf("number of installed packages: %v\n", len(installedPackages))

	pacman.Update(packages)
}

func readPackagesFile() pacman.Packages {
	var packages []pacman.Package

	pfile, err := ioutil.ReadFile("packages.json")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err = json.Unmarshal(pfile, &packages); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return packages
}

func main() {
	bootPtr := flag.Bool("boot", false, "setup live OS on target volume")
	installPtr := flag.Bool("install", false, "setup permanent OS on target volume")
	updatePtr := flag.Bool("update", false, "updates OS on target volume")

	flag.Parse()

	packages := readPackagesFile()

	if *bootPtr && !*installPtr && !*updatePtr {
		boot.Setup(packages)
		fmt.Println("when ready run: sudo dd if=archlinux-yyyy.mm.dd-x86_64.iso of=/dev/sdx status=progress")
	} else if !*bootPtr && *installPtr && !*updatePtr {
		if args := flag.Args(); len(args) == 0 {
			fmt.Println("install requires a target e.g sysconf -install /dev/sdx")
			os.Exit(1)
		}
		installOperation(flag.Args()[0])
	} else if !*bootPtr && !*installPtr && *updatePtr {
		update(packages)
	}

}
