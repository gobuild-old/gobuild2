package slave

import (
	"github.com/gobuild/log"
	"github.com/qiniu/api/io"
)

var defaultBulket = "xxxx"

type Storager interface {
	Upload(localFile string) (pubAddr string, err error)
}

type Qiniu struct {
	uptoken string
	key     string
	bulket  string
}

func (q *Qiniu) Upload(local string) (pubAddr string, err error) {
	var ret io.PutRet
	var extra = &io.PutExtra{}
	log.Infof("upload(qiniu) token: %s, key: %s, local: %s", q.uptoken, q.key, local)
	if err = io.PutFile(nil, &ret, q.uptoken, q.key, local, extra); err != nil {
		return
	}
	pubAddr = "http://" + q.bulket + ".qiniudn.com/" + q.key
	return
}
