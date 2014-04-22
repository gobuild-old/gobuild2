package routers

import "github.com/martini-contrib/render"

func Home(r render.Render) {
	r.HTML(200, "home", nil)
}
