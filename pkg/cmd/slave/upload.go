package slave

import (
	"github.com/qiniu/api/io"
	"github.com/qiniu/api/rs"
)

var Bulket = "gobuild-io"

func UploadQiniu(localFile string, destName string) (addr string, err error) {
	policy := rs.PutPolicy{Scope: Bulket}
	uptoken := policy.Token(nil)

	var ret io.PutRet
	var extra = new(io.PutExtra)
	if err = io.PutFile(nil, &ret, uptoken, destName, localFile, extra); err != nil {
		return
	}
	addr = "http://" + Bulket + ".qiniudn.com/" + destName
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
