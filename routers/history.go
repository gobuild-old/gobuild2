package routers

import (
	"net/http"
	"strconv"

	"github.com/go-martini/martini"
	"github.com/gobuild/gobuild2/models"
	"github.com/gobuild/middleware"
	"github.com/qiniu/log"
)

func History(ctx *middleware.Context, params martini.Params, req *http.Request) {
	id, _ := strconv.Atoi(req.FormValue("id"))
	tid := int64(id)
	task, err := models.GetTaskById(tid)
	if err != nil {
		log.Errorf("get task by id error: %v", err)
	}
	history, err := models.GetAllBuildHistoryByTid(tid)
	if err != nil {
		log.Errorf("get task history error: %v", err)
	}
	ctx.Data = map[string]interface{}{
		"Task":        task,
		"History":     history,
		"AutoRefresh": ctx.Query("auto_refresh") == "true",
	}
	ctx.HTML(200, "history")
}
