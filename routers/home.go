package routers

import (
	"strconv"

	"github.com/codegangsta/martini-contrib/web"
	"github.com/gobuild/gobuild2/models"
	"github.com/martini-contrib/render"
	"github.com/qiniu/log"
)

type RepoInfoForm struct {
	Name string `form:"repo-name" binding:"required"`
}

type RepositoryForm struct {
	Rid int64 `form:"rid" binding:"required"`
}

func NewRepo(rf RepoInfoForm, ctx *web.Context) {
	var err error
	if _, err = models.CreateRepository(rf.Name); err != nil {
		log.Errorf("create repo error: %v", err)
	}
	ctx.Redirect(302, "/")
}

func NewBuild(rf RepositoryForm, ctx *web.Context) {
	task := new(models.Task)
	task.Arch = "amd64"
	task.Os = "darwin"
	task.CommitId = "xxxxxx"
	task.CgoEnable = false
	task.Rid = rf.Rid
	if _, err := models.CreateTask(task); err != nil {
		log.Errorf("create module error: %v", err)
	}
	ctx.Redirect(302, "/repo?id="+strconv.Itoa(int(rf.Rid)))
}

func Home(r render.Render) {
	pv := models.RefreshPageView("/")
	repos, err := models.GetAllRepos(50, 0)
	if err != nil {
		log.Errorf("get repos from db error: %v", err)
	}
	r.HTML(200, "home", map[string]interface{}{
		"Repos": repos,
		"PV":    pv,
	})
}

func Ruok() string {
	return "imok"
}
