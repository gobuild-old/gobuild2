package routers

import (
	"github.com/gobuild/gobuild2/models"

	"github.com/gobuild/log"
	"github.com/gobuild/middleware"
)

type RepoInfoForm struct {
	Name string `form:"repo-name" binding:"required"`
}

type RepositoryForm struct {
	Rid int64 `form:"rid" binding:"required"`
}

type TaskForm struct {
	Tid int64 `form:"tid" binding:"required"`
}

func NewRepo(rf RepoInfoForm, ctx *middleware.Context) {
	ctx.Redirect(302, "/"+rf.Name)
}

func Home(ctx *middleware.Context) {
	pv := models.RefreshPageView("/")
	repos, err := models.GetAllRepos(50, 0)
	if err != nil {
		log.Errorf("get repos from db error: %v", err)
	}
	ctx.Data["Title"] = "home"
	ctx.Data["Repos"] = repos
	ctx.Data["PV"] = pv
	ctx.HTML(200, "home")
}

func Ruok() string {
	return "imok"
}

func NotFound(ctx *middleware.Context) {
	ctx.Handle(404, "Where you got this page", nil)
}
