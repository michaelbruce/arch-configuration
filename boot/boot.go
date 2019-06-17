package boot

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
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

	usr, err := user.Current()

	if err != nil {
		fmt.Println("could not find current user: ", err)
		os.Exit(1)
	}

	_, err = exec.Command("cp", "-rL", usr.HomeDir+"/.vim", dir+"/airootfs/root/.vim").Output()

	if err != nil {
		fmt.Println("could not copy .vim: ", err)
		os.Exit(1)
	}

	if err = os.Mkdir(dir+"/airootfs/root/.config", 0755); err != nil {
		fmt.Println("could not copy .vim: ", err)
		os.Exit(1)
	}

	_, err = exec.Command("cp",
		"-r",
		usr.HomeDir+"/.config/openbox",
		dir+"/airootfs/root/.config/openbox").Output()

	if err != nil {
		fmt.Println("could not copy .mozilla: ", err)
		os.Exit(1)
	}

	_, err = exec.Command("cp",
		"-rL",
		usr.HomeDir+"/.mozilla",
		dir+"/airootfs/root/.mozilla").Output()

	if err != nil {
		fmt.Println("could not copy .config/openbox: ", err)
		os.Exit(1)
	}

	err = copyFiles([]copyInstruction{
		copyInstruction{".bashrc", ".config/.bashrc"},
		copyInstruction{".xinitrc", ".xinitrc"},
		copyInstruction{".vimrc", ".vimrc"},
		copyInstruction{".Xresources", ".Xresources"},
		copyInstruction{".Xmodmap", ".Xmodmap"},
		copyInstruction{".xinitrc", ".xinitrc"},
		copyInstruction{".tmux.conf", ".tmux.conf"},
		copyInstruction{"src/sysconf/sysconf", "sysconf"},
		copyInstruction{"src/sysconf/packages.json", "packages.json"},
	}, dir)

	if err != nil {
		fmt.Println("problem copying files: %v", err)
	}

	appendSetBashInstructions(dir)
	build(dir)

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

type copyInstruction struct {
	Source      string
	Destination string
}

// copy a file from src to dst
func copy(src, dst string) error {
	srcStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !srcStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func copyFiles(instructions []copyInstruction, dir string) error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	for _, inst := range instructions {
		err := copy(usr.HomeDir+"/"+inst.Source, dir+"/airootfs/root/"+inst.Destination)
		if err != nil {
			return err
		}
	}

	return nil
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

func appendSetBashInstructions(dir string) error {
	f, err := os.OpenFile(dir+"/airootfs/root/customize_airootfs.sh", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString("usermod -s /usr/bin/bash root\n")

	if err != nil {
		return err
	}

	_, err = f.WriteString("cp $(which vim) /usr/local/bin/vi\n")

	if err != nil {
		return err
	}

	_, err = f.WriteString("mv /root/.config/.bashrc /root/.bashrc\n")
	return err
}

func build(dir string) {
	cmd := exec.Command("sudo", dir+"/build.sh", "-v", "-o", "./")
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	cmd.Wait()

	_, err := exec.Command("sudo", "rm", "-rf", "./work").Output()
	if err != nil {
		fmt.Println("could not remove work folder: ", err)
		os.Exit(1)
	}
}
