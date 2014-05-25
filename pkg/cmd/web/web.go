package web

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gobuild/gobuild2/models"
	"github.com/gobuild/gobuild2/pkg/base"
	"github.com/gobuild/gobuild2/pkg/config"
	"github.com/gobuild/gobuild2/pkg/xrpc"
	"github.com/gobuild/gobuild2/routers"
	"github.com/gobuild/log"

	"github.com/codegangsta/cli"
	"github.com/go-martini/martini"
	"github.com/gobuild/middleware"
	"github.com/martini-contrib/binding"
)

func newMartini() *martini.ClassicMartini {

	r := martini.NewRouter()
	m := martini.New()
	m.Use(Logger())
	m.Use(martini.Recovery())
	// m.Use(martini.Static("public"))
	m.Use(martini.Static("public", martini.StaticOptions{SkipLogging: true}))

	m.MapTo(r, (*martini.Routes)(nil))
	m.Action(r.Handle)

	var funcMap = base.TemplateFuncs

	m.Use(middleware.ContextWithCookieSecret("", middleware.Options{
		Funcs: []template.FuncMap{funcMap},
	}))
	return &martini.ClassicMartini{m, r}
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

	xrpc.HandleRpc()
	m.Get("/ruok", routers.Ruok)
	m.Any("/", routers.Home)
	m.Any("/repo", routers.Repo)
	m.Any("/history", routers.History)
	m.Any("/download", routers.Download)
	m.Post("/new-repo", binding.Bind(routers.RepoInfoForm{}), routers.NewRepo)
	m.Any("/search", routers.Search)

	m.Group("/api", func(r martini.Router) {
		m.Get("/pkglist", routers.PkgList)
		m.Post("/build", binding.Bind(routers.RepositoryForm{}), routers.NewBuild)
		m.Post("/force-rebuild", binding.Bind(routers.TaskForm{}), routers.ForceRebuild)
	})

	// Not found handler.
	m.NotFound(routers.NotFound)

	http.Handle("/", m)

	if err = models.ResetAllTaskStatus(); err != nil {
		log.Fatalf("reset all task status: %v", err)
	}
	go drainTask()

	listenAddr := fmt.Sprintf("%s:%d",
		cfg.Server.Addr,
		cfg.Server.Port)
	log.Printf("listen %s\n", strconv.Quote(listenAddr))
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
