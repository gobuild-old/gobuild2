package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

const VERSION = "0.0.1"

func runInit(c *cli.Context) {
	fmt.Println(c.Args())
	println("init")
}

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
			Name:   "init",
			Usage:  "initial gobuild.yml file",
			Action: runInit,
		},
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
