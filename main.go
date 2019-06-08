package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sysconf/pacman"
)

func install() {
	fmt.Println("installing Arch Linux on /dev/sdx...")
	// 1. Use lsblk -J to check that the target has at least 10gb
	// 2. Warn the user that you are going to wipe the target (or call them chicken.)
	// 3. determine if uefi or mbr machine with dmesg (dmesg | grep EFI >/dev/null)
	// TODO it would be nice to keep this setting - but perhaps we can just check whether either mbr or uefi packages are present and figure things out that way.
	// 4. check you are root/have the correct permissions
	// 5. setup the target disk reformat/create partitions
	// TODO: should I create swap partition or swap file (how do I make a swap file permanent?)
	// 6. include pacman with pacstrap
	// 7. create file system table
	// 8. download tools.. setup for root?
}

func update() {
	fmt.Println("updating system configuration...")

	var packages []pacman.Package

	pfile, err := ioutil.ReadFile("packages.json")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = json.Unmarshal(pfile, &packages)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if !pacman.Exists() {
		os.Exit(1)
	}

	fmt.Printf("packages to include %v\n", packages)

	installedPackages := pacman.InstalledPackages()

	fmt.Printf("number of installed packages: %v\n", len(installedPackages))

	pacman.Update(packages)
}

func main() {
	bootPtr := flag.Bool("boot", false, "setup live OS on target volume")
	installPtr := flag.Bool("install", false, "setup permanent OS on target volume")
	updatePtr := flag.Bool("update", false, "updates OS on target volume")

	flag.Parse()

	if *bootPtr && !*installPtr && !*updatePtr {
		fmt.Println("install boot OS")
	} else if !*bootPtr && *installPtr && !*updatePtr {
		install()
	} else if !*bootPtr && !*installPtr && *updatePtr {
		update()
	}

}
