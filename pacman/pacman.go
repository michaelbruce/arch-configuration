// Pacman is a package for interacting with Arch's Pacman Package Manager.
// https://www.archlinux.org/pacman/
package pacman

import (
	"fmt"
	"os/exec"
)

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
