package routers

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
)

func Register(m *martini.ClassicMartini) {
	m.Get("/ruok", Ruok)
	m.Any("/", Home)
	m.Any("/doc", Doc)
	m.Any("/repo", Repo)
	m.Any("/history", History)
	m.Any("/download", Download)
	m.Post("/new-repo", binding.Bind(RepoInfoForm{}), NewRepo)
	m.Any("/search", Search)

	m.Group("/api", func(r martini.Router) {
		m.Get("/pkglist", PkgList)
		m.Post("/build", binding.Bind(RepositoryForm{}), NewBuild)
		m.Post("/force-rebuild", binding.Bind(TaskForm{}), ForceRebuild)
	})

	m.Get("/**", Repo) // for the rest of request
	// Not found handler.
	// m.NotFound(routers.NotFound)
}
