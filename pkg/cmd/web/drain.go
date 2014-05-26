package web

import (
	"time"

	"github.com/gobuild/gobuild2/models"
	"github.com/gobuild/gobuild2/routers"
	"github.com/qiniu/log"
)

// keep task status always avaliable
func drainTask() {
	for {
		log.Infof("drain task start after 25min")
		time.Sleep(25 * time.Minute)
		if repos, err := models.GetAllRepos(1000, 0); err == nil {
			for _, r := range repos {
				routers.TriggerBuildRepositoryById(r.Id)
			}
		}
	}
}
