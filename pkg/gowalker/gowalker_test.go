package gowalker

import "testing"

var (
	pkgname     = "github.com/codeskyblue/fswatch"
	pkgjs       = "github.com/defunkt/dotjs"
	pkgnotexist = "github.com/xxxxx/aaddu98ka"
	pkglib      = "github.com/codeskyblue/go-sh"
)

func TestRefresh(t *testing.T) {
	err := RefreshPkg(pkgname)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetPkgInfo(t *testing.T) {
	pkginfo, err := GetPkgInfo(pkgname)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(pkginfo)
	pkginfo, err = GetPkgInfo(pkgjs)
	if err == nil {
		t.Fatal("should error here")
	}
	pkginfo, err = GetCmdPkgInfo(pkglib)
	if err == nil {
		t.Fatal("should error here")
	}
}
