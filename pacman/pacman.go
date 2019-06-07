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
// FOR NOW we are just including everything in requested packages
//var CorePackages = []string{
//	"ntp",
//}

// calls pacman and returns the result as a string
func pacman(args ...string) string {
	var cmd = exec.Command("pacman", args...)

	output, err := cmd.Output()

	if err != nil {
		fmt.Printf("pacman error: %v %v\n", args, err)
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

// uniq dedupes a list of Packages
func (ps Packages) uniq() Packages {
	seen := make(map[Package]struct{}, len(ps))
	j := 0
	for _, v := range ps {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		ps[j] = v
		j++
	}
	return ps[:j]
}

// Dependencies returns a list of packages that Packages ps depend on.
// We call dependencies as a group to perform a single remote query
// TODO fails when querying for a package that does not exist
// TODO does not support groups e.g "xorg" must use packages e.g "xorg-server"
func (ps Packages) Dependencies() Packages {
	var dependencies Packages
	var packageNames []string

	for _, p := range ps {
		packageNames = append(packageNames, p.Name)
	}

	args := append([]string{"-Si"}, packageNames...)
	var output = pacman(args...)

	for _, irow := range strings.Split(output, "\n") {
		if strings.Contains(irow, "Depends On") {
			for _, dep := range strings.Fields(irow)[3:] {
				dependencies = append(dependencies, Package{dep})
			}
		}
	}

	if len(dependencies) != 0 && dependencies[0].Name == "None" {
		return Packages{} // empty list when there are no dependencies
	}

	return dependencies
}

// Update takes a list of required packages, find their dependencies and installs them.
// Removes packages not included in this group except for core packages such as the Linux Kernel.
func Update(requested Packages) {
	fmt.Println("updating...")
	fmt.Printf("requested packages: %v\n", requested)

	var required = requested
	required = append(required, requested.Dependencies()...)

	fmt.Printf("required packages: %v\n", required)

	fmt.Printf("number of requested packages: %v\n", len(requested))
	fmt.Printf("number of required packages: %v\n", len(required))
	// list currently installed packages
	// find deps for requested packages
	// group that is not a requested, dep or core package named orphaned packages
	// remove orphans.
	// include requested + deps (group named required packages
}
