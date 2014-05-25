package web

import (
	"time"

	"github.com/gobuild/gobuild2/models"
	"github.com/qiniu/log"
)

// keep task status always avaliable
func drainTask() {
	for {
		log.Infof("drain task start")
		if repos, err := models.GetAllRepos(1000, 0); err == nil {
			for _, r := range repos {
				oas := map[string]string{
					"windows": "386",
					"linux":   "386",
					"darwin":  "amd64",
				}
				for os, arch := range oas {
					// err :=
					models.CreateNewBuilding(r.Id, "master", os, arch)
					// if err != nil {
					// log.Errorf("drain - create module error: %v", err)
					// }
				}
			}
		}
		time.Sleep(1 * time.Minute)
	}
}
