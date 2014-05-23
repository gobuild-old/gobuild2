package routers

import "github.com/gobuild/middleware"

type PackageItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func Search(ctx *middleware.Context) {
	var result []PackageItem
	result = append(result, PackageItem{
		Name:        "github.com/nsf/gocode",
		Description: "golang code complete",
	})
	ctx.JSON(200, result)
}
