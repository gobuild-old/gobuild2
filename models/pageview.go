package models

import (
	"sync"
	"time"

	"github.com/qiniu/log"
)

type PageView struct {
	Uri        string `xorm:"unique(pk)"`
	TotalCount int64
}

func init() {
	tables = append(tables, new(PageView))
}

var pageView = make(map[string]int64)

func drainPv() {
	for {
		time.Sleep(time.Second * 10)
		for uri, cnt := range pageView {
			affec, err := orm.Update(&PageView{TotalCount: cnt}, &PageView{Uri: uri})
			if err != nil || affec == 0 {
				orm.Insert(&PageView{Uri: uri, TotalCount: cnt})
			}
		}
	}
}

var pvOnce sync.Once

func RefreshPageView(uri string) int64 {
	pvOnce.Do(func() {
		orm.Insert(&PageView{uri, 0})
		go drainPv()
	})
	if pageView[uri] == 0 {
		pv := &PageView{Uri: uri}
		_, err := orm.Get(pv)
		if err != nil {
			log.Errorf("get pv from db error: %v", err)
			return -1
		}
		pageView[uri] = pv.TotalCount
	}
	pageView[uri] += 1
	return pageView[uri]
}
