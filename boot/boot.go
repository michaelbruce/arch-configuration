package boot

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

func Setup() {
	fmt.Println("create boot image")
	dir := config()
	defer os.RemoveAll(dir)
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
