package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gobuild/gobuild2/routers"
	"github.com/gobuild/log"

	"github.com/codegangsta/cli"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
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
func newMartini() *martini.ClassicMartini {
	/*
		r := martini.NewRouter()
		m := martini.New()
		m.Use(middleware.Logger())
		m.Use(martini.Recovery())
		m.Use(martini.Static("public"))
		m.MapTo(r, (*martini.Routes)(nil))
		m.Action(r.Handle)
	*/
	//return &martini.ClassicMartini{m, r}
	m := martini.Classic()
	m.Use(render.Renderer())
	return m
}

func initConfig(confPath string) *Config {
	cfg, err := readCfg(confPath)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}

func runWeb(c *cli.Context) {
	cfg := initConfig(c.String("conf"))
	m := newMartini()
	m.Get("/echo", func() string {
		return "hello"
	})
	m.Any("/", routers.Home)
	listenAddr := fmt.Sprintf("%s:%d",
		cfg.Server.Addr,
		cfg.Server.Port)
	log.Printf("listen %s\n", strconv.Quote(listenAddr))
	log.Fatal(http.ListenAndServe(listenAddr, m))
}
