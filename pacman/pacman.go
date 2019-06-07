// Pacman is a package for interacting with Arch's Pacman Package Manager.
// https://www.archlinux.org/pacman/
package pacman

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Package struct {
	Name string `json:"name"`
}

type Packages []Package

// depending sys config, may include uefi/mbr packages etc.
// might be better off as a function call to figure that out.
// CorePackages...

// calls pacman and returns the result as a string
func pacman(args ...string) string {
	var cmd = exec.Command("pacman", args...)

	output, err := cmd.Output()

	if err != nil {
		fmt.Printf("pacman not detected: %v\n", err)
		os.Exit(1)
	}

	return string(output)
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
func InstalledPackages() Packages {
	var output = pacman("-Q")
	var packages Packages

	packagesByLine := strings.Split(output, "\n")
	packagesByLine = packagesByLine[:len(packagesByLine)-1] // last line is always blank

	for _, pline := range packagesByLine {
		fields := strings.Fields(pline)
		packages = append(packages, Package{fields[0]})
	}

	return packages
}

// Dependencies returns a list of packages that Package p depends on.
func (p Package) Dependencies() Packages {
	var dependencies Packages
	var output = pacman("-Qi", p.Name)

	for _, irow := range strings.Split(output, "\n") {
		if strings.Contains(irow, "Depends On") {
			for _, dep := range strings.Fields(irow)[3:] {
				dependencies = append(dependencies, Package{dep})
			}
		}
	}

	return dependencies
}

// dedupes Packages
//func (ps Packages) Uniq() []Package {
//
//}

// Update takes a list of required packages, find their dependencies and installs them.
// Removes packages not included in this group except for core packages such as the Linux Kernel.
func Update(requested Packages) {
	fmt.Println("updating...")
	fmt.Printf("requested packages: %v\n", requested)

	var required Packages

	for _, p := range requested {
		required = append(required, p) // also want p.Dependencies...
		required = append(required, p.Dependencies()...)
	}

	fmt.Printf("required packages: %v\n", required)
	// list currently installed packages
	// find deps for requested packages
	// group that is not a requested, dep or core package named orphaned packages
	// remove orphans.
	// include requested + deps (group named required packages
}
