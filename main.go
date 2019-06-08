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
