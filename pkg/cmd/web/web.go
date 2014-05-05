package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gobuild/gobuild2/models"
	"github.com/gobuild/gobuild2/pkg/config"
	"github.com/gobuild/gobuild2/routers"
	"github.com/gobuild/log"

	"github.com/codegangsta/cli"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

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

func Action(c *cli.Context) {
	var err error
	if err = config.Load(c.String("conf")); err != nil {
		log.Fatal(err)
	}
	if err = models.InitDB(); err != nil {
		log.Fatal(err)
	}
	cfg := config.Config
	m := newMartini()
	m.Get("/ruok", routers.Ruok)
	m.Any("/", routers.Home)
	listenAddr := fmt.Sprintf("%s:%d",
		cfg.Server.Addr,
		cfg.Server.Port)
	log.Printf("listen %s\n", strconv.Quote(listenAddr))
	log.Fatal(http.ListenAndServe(listenAddr, m))
}
