package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gobuild/log"

	"github.com/codegangsta/cli"
	"github.com/go-martini/martini"
)

func init() {
	c := cli.Command{
		Name:   "web",
		Usage:  "start gobuild web server",
		Action: runWeb,
		Flags: []cli.Flag{
			cli.StringFlag{"conf,f", "conf/app.ini", "config file"},
		},
	}
	app.Commands = append(app.Commands, c)
}

func runWeb(c *cli.Context) {
	cfgPath := c.String("conf")
	cfg, err := readCfg(cfgPath)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(cfg)

	m := martini.Classic()
	m.Get("/", func() string {
		return "hello gobuild"
	})
	listenAddr := fmt.Sprintf("%s:%d",
		cfg.Server.Addr,
		cfg.Server.Port)
	log.Printf("listen %s\n", strconv.Quote(listenAddr))
	log.Fatal(http.ListenAndServe(listenAddr, m))
}
