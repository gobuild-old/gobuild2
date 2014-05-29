package models

import (
	"sync"
	"time"

	"github.com/gobuild/log"
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

func RefreshPageView(uri string, add ...int64) int64 {
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
	if len(add) == 0 {
		pageView[uri] += 1
	} else {
		pageView[uri] += add[0]
	}
	return pageView[uri]
}
