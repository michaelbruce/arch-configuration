package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
		if args := flag.Args(); len(args) == 0 {
			fmt.Println("install requires a target e.g sysconf -install /dev/sdx")
			os.Exit(1)
		}
		installOperation(flag.Args()[0])
	} else if !*bootPtr && !*installPtr && *updatePtr {
		update()
	}

}
