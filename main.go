package main

import (
	"log"

	"github.com/terrariumcloud/terrarium/cmd"
)

var buildInformationVersion string

func main() {
	cmd.Execute(buildInformationVersion)
	log.Println(buildInformationVersion)
}
