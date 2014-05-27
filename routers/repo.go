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
	models.RefreshPageView("/d/" + ctx.Query("rid"))
	ctx.Redirect(302, task.ArchieveAddr)
}

func TriggerBuildRepositoryById(rid int64) (err error) {
	repo, err := models.GetRepositoryById(rid)
	if err != nil {
		log.Errorf("get repo by id error: %v", err)
		return
	}
	oas := map[string]string{
		"windows": "386",
		"linux":   "386",
		"darwin":  "amd64",
	}
	if repo.IsCgo {
		delete(oas, "windows")
		oas["linux"] = "amd64"
	}
	for os, arch := range oas {
		err := models.CreateNewBuilding(rid, "master", os, arch)
		if err != nil {
			log.Errorf("create module error: %v", err)
		}
	}
	return nil
}

func NewBuild(rf RepositoryForm, ctx *middleware.Context) {
	defer ctx.Redirect(302, "/repo?id="+strconv.Itoa(int(rf.Rid)))
	TriggerBuildRepositoryById(rf.Rid)
}

func ForceRebuild(tf TaskForm, ctx *middleware.Context) {
	if err := models.ResetTask(tf.Tid); err != nil {
		log.Errorf("reset task failed: %v", err)
	}
	ctx.Redirect(302, "/history?id="+strconv.Itoa(int(tf.Tid))+"&auto_refresh=true")
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
		"DownCnt":    models.RefreshPageView("/d/"+ctx.Query("id"), 0),
	}
	rus, err := models.GetAllLastRepoUpdate(rid)
	if err != nil {
		log.Error("get last repo error: %v", err)
	}
	ctx.Data["Last"] = rus
	ctx.HTML(200, "repo")
}
