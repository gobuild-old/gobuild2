package routers

import (
	"net/http"
	"strconv"

	"github.com/go-martini/martini"
	"github.com/gobuild/gobuild2/models"
	"github.com/martini-contrib/render"
	"github.com/qiniu/log"
)

func Repo(r render.Render, params martini.Params, req *http.Request) {
	id, _ := strconv.Atoi(req.FormValue("id"))
	repo, err := models.GetRepositoryById(int64(id))
	if err != nil {
		log.Errorf("get single repo from db error: %v", err)
	}
	recentTask, _ := models.GetTaskById(1)
	r.HTML(200, "repo", map[string]interface{}{
		"Repo":       repo,
		"RecentTask": recentTask,
		"Tasks":      []*models.Task{recentTask},
	})
}
