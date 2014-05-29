package routers

import (
	"github.com/gobuild/gobuild2/models"
	"github.com/gobuild/log"
	"github.com/gobuild/middleware"
)

func Search(ctx *middleware.Context) {
	log.Info(ctx.Request.RequestURI)
	pv := models.RefreshPageView(ctx.Request.RequestURI) // "/search")
	repos, err := models.GetAllRepos(50, 0)
	if err != nil {
		log.Errorf("get repos from db error: %v", err)
	}
	ctx.Data["Title"] = "home"
	ctx.Data["Repos"] = repos
	ctx.Data["PV"] = pv
	ctx.HTML(200, "search")
}
