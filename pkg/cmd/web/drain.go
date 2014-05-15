package web

import "github.com/gobuild/gobuild2/models"

// keep task status always avaliable
func drainTask() {
	task := &models.Task{
		// Repo:      "github.com/codeskyblue/fswatch",
		Branch:    "master",
		CgoEnable: false,
		Arch:      "386",
		Os:        "windows",
	}
	models.CreateTask(task)
}
