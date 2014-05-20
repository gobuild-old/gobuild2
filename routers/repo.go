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
	rid := int64(id)
	repo, err := models.GetRepositoryById(rid)
	if err != nil {
		log.Errorf("get single repo from db error: %v", err)
	}
	tasks, err := models.GetTasksByRid(rid)
	if err != nil {
		log.Errorf("get tasks by id, error: %v", err)
	}
	recentTask, _ := models.GetTaskById(1)
	r.HTML(200, "repo", map[string]interface{}{
		"Repo":       repo,
		"RecentTask": recentTask,
		"Tasks":      tasks, // []*models.Task{recentTask, recentTask},
	})
}
