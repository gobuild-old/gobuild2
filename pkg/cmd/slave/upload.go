package slave

import (
	"strings"

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

func UploadQiniu(localFile string, destName string) (addr string, err error) {
	policy := rs.PutPolicy{Scope: defaultBulket}
	uptoken := policy.Token(nil)

	destName = strings.TrimLeft(destName, "/")
	var ret io.PutRet
	var extra = new(io.PutExtra)
	if err = io.PutFile(nil, &ret, uptoken, destName, localFile, extra); err != nil {
		return
	}
	addr = "http://" + defaultBulket + ".qiniudn.com/" + destName
	return
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
