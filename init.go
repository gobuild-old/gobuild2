package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
	_ = os.Stdout
}

const RCFILE = ".gobuild.yml"

type PackageConfig struct {
	Filesets struct {
		Includes []string `yaml:"includes"`
		Excludes []string `yaml:"excludes"`
	} `yaml:"filesets"`
	Settings struct {
		GoFlags   string `yaml:"goflags"`
		CGOEnable bool   `yaml"cgoenable"`
	}
}

var defaultDcfg *PackageConfig

func init() {
	pcfg := &PackageConfig{}
	pcfg.Filesets.Includes = []string{"README.md", "LICENSE"}
	pcfg.Filesets.Excludes = []string{".*.go"}
	pcfg.Settings.CGOEnable = true
	pcfg.Settings.GoFlags = ""
	defaultDcfg = pcfg
}

func runInit(c *cli.Context) {
	data, _ := goyaml.Marshal(defaultDcfg)
	if err := ioutil.WriteFile(RCFILE, data, 0644); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("rcfile: %s created", RCFILE)
	}
}
