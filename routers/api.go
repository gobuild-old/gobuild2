package routers

import (
	"time"

	"github.com/gobuild/middleware"
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
	var result []PackageItem
	result = append(result, PackageItem{
		Name:        "github.com/nsf/gocode",
		Description: "golang code complete",
		Branches:    []Branch{Branch{"master", "abcdeftg", fmtTime(time.Now())}},
		// Updated:     fmtTime(time.Now()),
	})
	ctx.JSON(200, result)
}
