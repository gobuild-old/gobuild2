package slave

import (
	"strings"

	"github.com/gobuild/log"
	"github.com/qiniu/api/conf"
	"github.com/qiniu/api/io"
	"github.com/qiniu/api/rs"
)

var defaultBulket = "xxxx"

func initQiniu(access, secret string, bulket string) {
	conf.ACCESS_KEY = access
	conf.SECRET_KEY = secret
	defaultBulket = bulket
}

type Storager interface {
	Upload(localFile string) (pubAddr string, err error)
}

type Qiniu struct {
	uptoken string
	key     string
}

func (q *Qiniu) Upload(local string) (pubAddr string, err error) {
	var ret io.PutRet
	var extra = &io.PutExtra{}
	if err = io.PutFile(nil, &ret, q.uptoken, q.key, local, extra); err != nil {
		return
	}
	log.Infof("upload success:%v", ret)
	pubAddr = "http://" + defaultBulket + ".qiniudn.com/" + q.key
	return
}

// mimetype ref: http://webdesign.about.com/od/multimedia/a/mime-types-by-content-type.htm
func UploadQiniu(localFile string, destName string) (addr string, err error) {
	key := strings.TrimLeft(destName, "/")
	policy := rs.PutPolicy{Scope: defaultBulket + ":" + destName}
	uptoken := policy.Token(nil)

	q := &Qiniu{uptoken, key}
	return q.Upload(localFile)
	// var ret io.PutRet
	// mimeType := ""
	// if strings.HasSuffix(destName, "tar.gz") {
	// 	mimeType = "application/x-tgz"
	// } else if strings.HasSuffix(destName, ".zip") {
	// 	mimeType = "application/zip"
	// }
	// var extra = &io.PutExtra{
	// 	MimeType: mimeType,
	// }
	// var extra = &io.PutExtra{}
	// if err = io.PutFile(nil, &ret, uptoken, destName, localFile, extra); err != nil {
	// 	return
	// }
	// log.Infof("upload success:%v", ret)
	// addr = "http://" + defaultBulket + ".qiniudn.com/" + destName
	// return
}

/*
func UploadLocal(file string) (addr string, err error) {
	f, err := ioutil.TempFile("files/", "upload-")
	if err != nil {
		return
	}
	err = f.Close()
	if err != nil {
		return
	}
	exec.Command("cp", "-f", file, f.Name()).Run()
	addr = "http://" + filepath.Join(Hostname, f.Name())
	return
}
*/
