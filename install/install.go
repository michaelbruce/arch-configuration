package install

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strconv"
)

type blockResponse struct {
	BlockDevices []blockInfo `json:"blockdevices"`
}

type blockInfo struct {
	Name string `json:"name"`
	Size string `json:"name"`
}

func CheckCapacity(name string) (bool, error) {
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
		fmt.Println(bi.Name, bi.Size)
		if "/dev/"+bi.Name == name {
			size, err := strconv.ParseFloat(bi.Size[:len(bi.Size)-1], 64)

			if err != nil {
				return false, err
			}

			if size > 10 {
				return true, nil
			} else {
				return false, fmt.Errorf("not enough space on device %v\n", name)
			}
		}
	}

	return false, fmt.Errorf("could not find block device %v\n", name)
}
