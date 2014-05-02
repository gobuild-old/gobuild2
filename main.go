package main

import (
	"os"

	"github.com/codegangsta/cli"
)

const VERSION = "0.0.1"

func runSlave(c *cli.Context) {
	println("slave")
}

var app = cli.NewApp()

func init() {
	app.Name = "gobuild"
	app.Usage = "[COMMANDS]"
	app.Version = VERSION
	app.Commands = append(app.Commands,
		cli.Command{
			Name:   "slave",
			Usage:  "start gobuild compile slave",
			Action: runSlave,
		},
		cli.Command{
			Name:   "pack",
			Usage:  "build and pack file into tgz or zip",
			Action: runPack,
		},
	)
}

func main() {
	app.Run(os.Args)
}
