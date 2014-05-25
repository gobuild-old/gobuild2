package routers

import (
	"time"

	"github.com/gobuild/gobuild2/models"
	"github.com/gobuild/middleware"
	"github.com/qiniu/log"
)

type PackageItem struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Branches    []Branch `json:"branches"`
}

type Branch struct {
	Name    string `json:"name"`
	Sha     string `json:"sha"`
	Updated string `json:"updated"`
	Os      string `json:"os"`
	Arch    string `json:"arch"`
}

func fmtTime(t time.Time) string { return t.UTC().Format(time.RFC3339) }

func PkgList(ctx *middleware.Context) {
	rs, err := models.GetAllLastRepoByOsArch(ctx.Query("os"), ctx.Query("arch"))
	if err != nil {
		ctx.JSON(400, nil)
		return
	}
	var result []*PackageItem
	var lastRid int64 = 0
	for _, lr := range rs {
		if lr.Rid == lastRid {
			if pos := len(result) - 1; pos >= 0 {
				r := result[pos]
				r.Branches = append(r.Branches,
					Branch{lr.TagBranch, lr.Sha, fmtTime(lr.Updated), lr.Os, lr.Arch})
			}
			continue
		}
		lastRid = lr.Rid
		repo, err := models.GetRepositoryById(lr.Rid)
		if err != nil {
			log.Errorf("a missing repo in last_repo_update: %v", lr)
			continue
		}
		br := Branch{lr.TagBranch, lr.Sha, fmtTime(lr.Updated), lr.Os, lr.Arch}
		result = append(result, &PackageItem{
			Name:        repo.Uri,   // "github.com/nsf/gocode",
			Description: repo.Brief, // "golang code complete",
			Branches:    []Branch{br},
		})

	}
	ctx.JSON(200, result)
}
