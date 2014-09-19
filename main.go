package main

import (
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/gobuild/gobuild2/cmd/pack"
	"github.com/gobuild/gobuild2/cmd/runinit"
	"github.com/gobuild/gobuild2/cmd/slave"
	"github.com/gobuild/gobuild2/cmd/web"
	"github.com/gobuild/log"
)

const VERSION = "0.0.1.0607"

var app = cli.NewApp()

func init() {
	cwd, _ := os.Getwd()
	program := filepath.Base(cwd)

	app.Name = "gobuild"
	app.Usage = "[COMMANDS]"
	app.Version = VERSION
	app.Commands = append(app.Commands,
		cli.Command{
			Name:   "slave",
			Usage:  "start gobuild compile slave",
			Action: slave.Action,
			Flags: []cli.Flag{
				cli.StringFlag{Name: "webaddr,w", Value: "localhost:8010", Usage: "gobuild2 web address"},
			},
		},
		cli.Command{
			Name:   "init",
			Usage:  "initial gobuild.yml file",
			Action: runinit.Action,
		},
		cli.Command{
			Name:   "pack",
			Usage:  "build and pack file into tgz or zip",
			Action: pack.Action,
			Flags: []cli.Flag{
				cli.StringFlag{Name: "os", Value: os.Getenv("GOOS"), Usage: "operation system"},
				cli.StringFlag{Name: "arch", Value: os.Getenv("GOARCH")},
				cli.StringFlag{Name: "depth", Value: "3", Usage: "depth of file to walk"},
				cli.StringFlag{Name: "output,o", Value: program + ".zip", Usage: "target file"},
				cli.StringFlag{Name: "gom", Value: "go", Usage: "go package manage program"},
				cli.BoolFlag{Name: "nobuild", Usage: "donot call go build when pack"},
				cli.StringSliceFlag{Name: "add,a", Value: &cli.StringSlice{}, Usage: "add file"},
			},
		},
		cli.Command{
			Name:   "web",
			Usage:  "start gobuild web server",
			Action: web.Action,
			Flags: []cli.Flag{
				cli.StringFlag{Name: "conf,f", Value: "conf/app.ini"},
			},
		},
	)
}

func main() {
	log.SetOutputLevel(log.Ldebug)
	app.Run(os.Args)
}
