package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sysconf/pacman"
)

// !/bin/sh -e
// keeps arch install fresh by:
//   - updating packages

func main() {
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

	// from the installed packages
	// what depends on my 'core' packages (from packages.json)
	// get all deps for 'core' packages
	// remove orphaned packages (handle cases like linux and linux-headers
	// define core packages by sysconf with additional 'required' packages from packages.json
	// install all 'core' and deps for core.
}
