// Pacman is a package for interacting with Arch's Pacman Package Manager.
// https://www.archlinux.org/pacman/
package pacman

import (
	"fmt"
	"os/exec"
	"strings"
)

type Package struct {
	Name string `json:"name"`
}

// Exists is used to make sure the pacman command exists.
func Exists() bool {
	var cmd = exec.Command("pacman", "--version")

	output, err := cmd.Output()

	if err != nil {
		fmt.Printf("pacman not detected: %v\n", err)
		return false
	}

	fmt.Println(string(output))

	return true
}

// Returns a list of all packages installed by pacman.
func InstalledPackages() []Package {
	var cmd = exec.Command("pacman", "-Q")

	output, err := cmd.Output()

	if err != nil {
		fmt.Printf("request for packages from pacman failed:", err)
		return []Package{}
	}

	var packages []Package

	packagesByLine := strings.Split(string(output), "\n")
	packagesByLine = packagesByLine[:len(packagesByLine)-1] // last line is always blank

	for _, pline := range packagesByLine {
		fields := strings.Fields(pline)
		packages = append(packages, Package{fields[0]})
	}

	return packages
}

// Update takes a list of required packages, find their dependencies and installs them.
// Removes packages not included in this group except for core packages such as the Linux Kernel.
func Update(required []Package) {
	fmt.Println("updating...")
	fmt.Printf("required packages: %v\n", required)
	// list currently installed packages
	// find deps for required packages
	// group that is not a required, dep or core package named orphaned packages
	// remove orphans.
	// include required + deps
}
