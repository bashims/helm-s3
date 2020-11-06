package main

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"os"
)

func main() {
	chartSemVer, err := semver.NewVersion("0.3.0-dev.0.0+20200926T155517Z")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(chartSemVer.String())
}
