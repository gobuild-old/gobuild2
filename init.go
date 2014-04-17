package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"launchpad.net/goyaml"

	"github.com/codegangsta/cli"
)

func init() {
	c := cli.Command{
		Name:   "init",
		Usage:  "initial gobuild.yml file",
		Action: runInit,
	}
	app.Commands = append(app.Commands, c)
}

const RCFILE = ".gobuild.yml"

type PackageConfig struct {
	Filesets struct {
		Includes []string `yaml:"includes"`
		Excludes []string `yaml:"excludes"`
	} `yaml:"filesets"`
}

func runInit(c *cli.Context) {
	pcfg := &PackageConfig{}
	pcfg.Filesets.Includes = []string{"README.md"}
	pcfg.Filesets.Excludes = []string{".*.go"}
	data, _ := goyaml.Marshal(pcfg)
	if err := ioutil.WriteFile(RCFILE, data, 0644); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("rcfile: %s created", RCFILE)
	}
}
