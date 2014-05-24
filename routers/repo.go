package routers

import (
	"net/http"
	"strconv"

	"github.com/go-martini/martini"
	"github.com/gobuild/gobuild2/models"
	"github.com/gobuild/middleware"
	"github.com/qiniu/log"
)

func Download(ctx *middleware.Context) {
	rid, _ := strconv.Atoi(ctx.Request.FormValue("rid"))
	os := ctx.Request.FormValue("os")
	arch := ctx.Request.FormValue("arch")
	task, err := models.GetOneDownloadableTask(int64(rid), os, arch)
	if err != nil {
		log.Errorf("get download task: %v", err)
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusNotFound)
		return
	}
	ctx.Redirect(302, task.ArchieveAddr)
}

func Repo(ctx *middleware.Context, params martini.Params, req *http.Request) {
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
	ctx.Data = map[string]interface{}{
		"Repo":       repo,
		"RecentTask": recentTask,
		"Tasks":      tasks,
	}
	ctx.HTML(200, "repo")
}
