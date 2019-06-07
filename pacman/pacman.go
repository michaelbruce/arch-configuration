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
		fmt.Printf("pacman error: %v %v %v\n", args, string(output), err)
		os.Exit(1)
	}

	return string(output)
}

func pactree(args ...string) string {
	var cmd = exec.Command("pactree", args...)

	output, err := cmd.Output()

	if err != nil {
		fmt.Printf("pactree error: %v %v %v\n", args, string(output), err)
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
// TODO could be querying local db if package exists
// TODO could be querying concurrently
func (ps Packages) Dependencies() Packages {
	var dependencies Packages
	var packageNames []string

	for _, p := range ps {
		packageNames = strings.Split(pactree("-slu", p.Name), "\n")
		packageNames = packageNames[:len(packageNames)-1]

		for _, d := range packageNames {
			dependencies = append(dependencies, Package{d})
		}
	}

	return dependencies.uniq()
}

// returns Packages that are missing from ps according to cps
func (ps Packages) missing(cps Packages) Packages {
	var missing Packages

	for _, cp := range cps {
		var found bool
		for _, p := range ps {
			if p.Name == cp.Name {
				found = true
			}
		}

		if !found {
			missing = append(missing, cp)
		}
	}

	return missing
}

func Remove(ps Packages) error {
	var names []string

	for _, p := range ps {
		names = append(names, p.Name)
	}

	args := []string{"pacman", "-R", "--noconfirm"}
	args = append(args, names...)
	cmd := exec.Command("sudo", args...)

	out, err := cmd.Output()

	if err != nil {
		fmt.Println(string(out), err)
	}

	fmt.Println(string(out))
	return nil
}

func Install(ps Packages) error {
	var names []string

	for _, p := range ps {
		names = append(names, p.Name)
	}

	args := []string{"pacman", "-Syu", "--noconfirm", "--needed"}
	args = append(args, names...)
	cmd := exec.Command("sudo", args...)

	out, err := cmd.Output()

	if err != nil {
		fmt.Println(string(out), err)
	}

	fmt.Println(string(out))
	return nil
}

// Update takes a list of required packages, find their dependencies and installs them.
// Removes packages not included in this group except for core packages such as the Linux Kernel.
func Update(requested Packages) {
	fmt.Println("updating...")
	fmt.Printf("requested packages: %v\n", requested)

	var required = requested
	required = append(required, requested.Dependencies()...)

	fmt.Printf("required packages: %v\n", required)

	extra := required.missing(InstalledPackages())
	fmt.Printf("extra packages: %v\n", extra)

	fmt.Printf("number of requested packages: %v\n", len(requested))
	fmt.Printf("number of required packages: %v\n", len(required))
	fmt.Printf("number of installed packages: %v\n", len(InstalledPackages()))
	fmt.Printf("number of extra packages: %v\n", len(extra))
	// remove extra packages
	// install required packages

	if len(extra) != 0 {
		fmt.Printf("Removing %v extra packages\n", len(extra))
		err := Remove(extra)

		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("Installing %v required packages\n", len(required))
	err := Install(required)

	if err != nil {
		panic(err)
	}
}
