package main

import (
	"github.com/codegangsta/cli"
	sh "github.com/codeskyblue/go-sh"
)

func runPack(c *cli.Context) {
	println("pack")
	sess := sh.NewSession()
	sess.ShowCMD = true
	sess.Command("gopm", "build").Run()
}
