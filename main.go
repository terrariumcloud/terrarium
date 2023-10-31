package main

import (
	"log"

	"github.com/terrariumcloud/terrarium/cmd"
)

var buildInformationVersion string

func main() {
	log.Println(buildInformationVersion)
	cmd.Execute(buildInformationVersion)
}
