package routers

import (
	"time"

	"github.com/gobuild/gobuild2/models"
	"github.com/gobuild/middleware"
	"github.com/qiniu/log"
)

type PackageItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	// Updated     string `json:"updated"`
	Branches []Branch `json:"branches"`
}

type Branch struct {
	Name    string `json:"name"`
	Sha     string `json:"sha"`
	Updated string `json:"updated"`
}

func fmtTime(t time.Time) string { return t.UTC().Format(time.RFC3339) }

func PkgList(ctx *middleware.Context) {
	rs, err := models.GetAllLastRepoByOsArch(ctx.Query("os"), ctx.Query("arch"))
	if err != nil {
		ctx.JSON(400, nil)
		return
	}
	var result []PackageItem
	for _, lr := range rs {
		repo, err := models.GetRepositoryById(lr.Rid)
		if err != nil {
			log.Errorf("a missing repo in last_repo_update: %v", lr)
			continue
		}
		result = append(result, PackageItem{
			Name:        repo.Uri,                                                    // "github.com/nsf/gocode",
			Description: repo.Brief,                                                  // "golang code complete",
			Branches:    []Branch{Branch{lr.TagBranch, lr.Sha, fmtTime(lr.Updated)}}, // "master", "abcdeftg", fmtTime(time.Now())}},
			// Updated:     fmtTime(time.Now()),
		})

	}
	ctx.JSON(200, result)
}
