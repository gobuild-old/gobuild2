package web

import (
	"github.com/gobuild/gobuild2/models"
	"github.com/qiniu/log"
)

func testAdd() {
	repo, err := models.CreateRepository("github.com/codeskyblue/fswatch")
	if err != nil {
		log.Errorf("create repo err: %v", err)
	}
	task := &models.Task{
		// Repo:      "github.com/codeskyblue/fswatch",
		Branch:    "master",
		CgoEnable: false,
		Arch:      "386",
		Os:        "windows",
		Rid:       repo.Id,
	}
	models.CreateTask(task)
}

// keep task status always avaliable
func drainTask() {
	// models.GetAvaliableTask(true, "windows", "386")
}
