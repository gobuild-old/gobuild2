package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"launchpad.net/goyaml"

	"github.com/codegangsta/cli"
	"github.com/gobuild/gobuild2/pkg/config"
)

func runInit(c *cli.Context) {
	data, _ := goyaml.Marshal(config.DefaultPcfg)
	if err := ioutil.WriteFile(config.RCFILE, data, 0644); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("rcfile: %s created", config.RCFILE)
	}
}
