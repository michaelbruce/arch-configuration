package install

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type blockResponse struct {
	BlockDevices []blockInfo `json:"blockdevices"`
}

type blockInfo struct {
	Name string `json:"name"`
	Size string `json:"size"`
}

// Ensure the target volume is big enough to comfortably contain Arch GNU/Linux
func CheckCapacity(name string) error {
	cmd := exec.Command("lsblk", "-J")

	out, err := cmd.Output()

	if err != nil {
		log.Fatal(err)
	}

	var response blockResponse

	err = json.Unmarshal(out, &response)

	if err != nil {
		log.Fatal(err)
	}

	for _, bi := range response.BlockDevices {
		fmt.Println("size is", bi.Size)
		if "/dev/"+bi.Name == name {
			size, err := strconv.ParseFloat(bi.Size[:len(bi.Size)-1], 64)

			if err != nil {
				return err
			}

			if size > 10 {
				return nil
			} else {
				return fmt.Errorf("not enough space on device %v\n", name)
			}
		}
	}

	return fmt.Errorf("could not find block device %v\n", name)
}

// Ensure the target machine has a UEFI
func CheckUEFI() error {
	cmd := exec.Command("dmesg")

	out, err := cmd.Output()

	if err != nil {
		log.Fatal(err)
	}

	kernelRingBuffer := strings.Split(string(out),"\n")
	kernelRingBuffer = kernelRingBuffer[:len(kernelRingBuffer)-1]

	for _, line := range kernelRingBuffer {
		if strings.Contains(line, "UEFI") {
			return nil
		}
	}

	return fmt.Errorf("no entry in the kernel ring buffer refers to UEFI")
}
