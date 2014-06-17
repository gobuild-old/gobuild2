package routers

import (
	"net/http"
	"strconv"

	"github.com/go-martini/martini"
	"github.com/gobuild/gobuild2/models"
	"github.com/gobuild/gobuild2/pkg/base"
	"github.com/gobuild/gobuild2/pkg/config"
	"github.com/gobuild/log"
	"github.com/gobuild/middleware"
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
	ctx.Redirect(302, task.ZipBallUrl)
}

func TriggerBuildRepositoryById(rid int64) (err error) {
	repo, err := models.GetRepositoryById(rid)
	if err != nil {
		log.Errorf("get repo by id error: %v", err)
		return
	}
	cvsinfo, _ := base.ParseCvsURI(repo.Uri)
	defaultBranch := "master"
	if cvsinfo.Provider == base.PVD_GOOGLE {
		defaultBranch = ""
	}
	models.CreateNewBuilding(rid, defaultBranch, "", "", models.AC_SRCPKG)
	if !repo.IsCmd {
		return nil
	}
	oas := map[string]string{
		"windows": "386",
		"windows": "amd64",
		"linux":   "386",
		"linux":   "amd64",
		"darwin":  "amd64",
	}
	if repo.IsCgo {
		delete(oas, "windows")
		oas["linux"] = "amd64"
	}
	for os, arch := range oas {
		err := models.CreateNewBuilding(rid, defaultBranch, os, arch, models.AC_BUILD)
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
	reponame := params["_1"]
	var repo *models.Repository
	var err error
	var tasks []models.Task
	var recentTask *models.Task
	if reponame != "" {
		repo, err = models.GetRepositoryByName(reponame)
		if err == models.ErrRepositoryNotExists {
			r, er := models.AddRepository(reponame)
			if er != nil {
				err = er
				ctx.Data["Error"] = err.Error()
				ctx.HTML(200, "repo")
				return
			}
			TriggerBuildRepositoryById(r.Id)
			ctx.Redirect(302, "/"+r.Uri)
			return
		}
		if err != nil {
			log.Errorf("get single repo from db error: %v", err)
		}
	} else {
		id, _ := strconv.Atoi(req.FormValue("id"))
		rid := int64(id)
		repo, err = models.GetRepositoryById(rid)
		if err != nil {
			log.Errorf("get single repo from db error: %v", err)
		}
	}
	tasks, err = models.GetTasksByRid(repo.Id)
	if err != nil {
		log.Errorf("get tasks by id, error: %v", err)
	}
	recentTask, _ = models.GetTaskById(1)
	ctx.Data["Repo"] = repo
	ctx.Data["RecentTask"] = recentTask
	ctx.Data["Tasks"] = tasks
	ctx.Data["DownCnt"] = models.RefreshPageView("/d/"+base.ToStr(repo.Id), 0)
	ctx.Data["RootUrl"] = config.Config.Server.RootUrl
	rus, err := models.GetAllLastRepoUpdate(repo.Id)
	if err != nil {
		log.Error("get last repo error: %v", err)
	}
	ctx.Data["Last"] = rus
	ctx.HTML(200, "repo")
}
