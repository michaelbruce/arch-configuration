package boot

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sysconf/pacman"
)

func Setup(packages pacman.Packages) {
	fmt.Println("create boot image")

	var bootPackages pacman.Packages
	for _, p := range packages {
		if p.Type == "boot" {
			bootPackages = append(bootPackages, p)
		}
	}

	dir := config()
	err := appendBootPackages(dir, bootPackages)

	if err != nil {
		fmt.Println("failed to append boot packages: ", err)
		os.Exit(1)
	}

	// defer os.RemoveAll(dir)
}

// copies across the release engineering config
// included by the archiso package
// returns the temporary path
func config() string {
	dir, err := ioutil.TempDir("", "sysconf-boot-")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("created directory", dir)

	cmd := exec.Command("cp", "-r", "/usr/share/archiso/configs/releng/.", dir)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()

	if err != nil {
		fmt.Println(string(out), stderr.String(), err)
		os.Exit(1)
	}

	return dir
}

// The packages.x86_64 file in the archiso defines which packages are installed
// on the live system. This function appends a few helpful extras.
func appendBootPackages(dir string, packages pacman.Packages) error {
	f, err := os.OpenFile(dir+"/packages.x86_64", os.O_APPEND|os.O_WRONLY, 0600)

	if err != nil {
		return err
	}

	defer f.Close()

	for _, p := range packages {
		if _, err = f.WriteString(p.Name + "\n"); err != nil {
			return err
		}
	}

	return nil
}
