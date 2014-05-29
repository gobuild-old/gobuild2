package routers

import (
	"strings"

	"github.com/gobuild/gobuild2/models"
	"github.com/gobuild/gobuild2/pkg/base"
	"github.com/gobuild/gobuild2/pkg/gowalker"
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
	defer ctx.Redirect(302, "/")
	var err error
	cvsinfo, err := base.ParseCvsURI(rf.Name) // base.SanitizedRepoPath(rf.Name)
	if err != nil {
		log.Errorf("parse cvs url error: %v", err)
		return
	}

	repoUri := cvsinfo.FullPath
	r := new(models.Repository)
	r.Uri = repoUri

	pkginfo, err := gowalker.GetCmdPkgInfo(repoUri)
	if err != nil {
		log.Errorf("gowalker not passed check: %v", err)
		return
	}
	r.IsCgo = pkginfo.IsCgo
	// description
	r.Brief = pkginfo.Description
	base.ParseCvsURI(repoUri)
	if strings.HasPrefix(repoUri, "github.com") {
		// comunicate with github
		fields := strings.Split(repoUri, "/")
		owner, repoName := fields[1], fields[2]
		repo, _, err := models.GHClient.Repositories.Get(owner, repoName)
		if err != nil {
			log.Errorf("get information from github error: %v", err)
		} else {
			r.Brief = *repo.Description
		}
	}
	if _, err = models.CreateRepository(r); err != nil {
		log.Errorf("create repo error: %v", err)
		return
	}
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
