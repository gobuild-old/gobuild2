package routers

import (
	"github.com/gobuild/gobuild2/models"
	"github.com/martini-contrib/render"
	"github.com/qiniu/log"
)

func Home(r render.Render) {
	repos, err := models.GetAllRepos(50, 0)
	if err != nil {
		log.Errorf("get repos from db error: %v", err)
	}
	r.HTML(200, "home", map[string]interface{}{
		"Repos": repos,
	})
}

func Ruok() string {
	return "imok"
}
