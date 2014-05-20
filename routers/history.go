package routers

import (
	"net/http"
	"strconv"

	"github.com/go-martini/martini"
	"github.com/gobuild/gobuild2/models"
	"github.com/martini-contrib/render"
	"github.com/qiniu/log"
)

func History(r render.Render, params martini.Params, req *http.Request) {
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
	r.HTML(200, "history", map[string]interface{}{
		"Task":    task,
		"History": history,
	})
}
