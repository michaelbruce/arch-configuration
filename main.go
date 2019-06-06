package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// !/bin/sh -e
// keeps arch install fresh by:
//   - updating packages

type Package struct {
	Name string `json:"name"`
}

func main() {
	fmt.Println("updating system configuration...")

	var packages []Package

	pfile, err := ioutil.ReadFile("packages.json")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = json.Unmarshal(pfile, &packages)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("packages %v\n", packages)
}
