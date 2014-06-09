package routers

import (
	"html/template"
	"io/ioutil"

	"github.com/gobuild/middleware"
	"github.com/russross/blackfriday"
)

func Doc(ctx *middleware.Context) {
	readme, err := ioutil.ReadFile("README.md")
	if err != nil {
		ctx.Data["Error"] = err.Error()
	}
	ctx.Data["Readme"] = template.HTML(string(blackfriday.MarkdownCommon(readme)))
	ctx.HTML(200, "doc")
}
