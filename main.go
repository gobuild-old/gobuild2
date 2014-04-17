package main

import (
	"os"

	"github.com/codegangsta/cli"
)

const VERSION = "0.0.1"

func runDaemon(c *cli.Context) {
	println("daemon")
}

var app = cli.NewApp()

func init() {
	app.Name = "gobuild"
	app.Usage = "[COMMANDS]"
	app.Version = VERSION
	app.Commands = append(app.Commands,
		cli.Command{
			Name:   "daemon",
			Usage:  "start gobuild compile daemon",
			Action: runDaemon,
		},
	)
}

func main() {
	app.Run(os.Args)
}
