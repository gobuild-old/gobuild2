package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/gobuild/gobuild2/modules/worker"
)

const VERSION = "0.0.1"

func runSlave(c *cli.Context) {
	println("slave")
	worker.Work()
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
	)
}

func main() {
	app.Run(os.Args)
}
