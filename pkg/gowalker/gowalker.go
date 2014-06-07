package gowalker

import (
	"crypto/tls"
	"errors"
	"net/http"

	"github.com/Unknwon/com"
)

// refresh: /refresh? pkgname=

var (
	base       = "https://gowalker.org/api/v1/"
	searchApi  = base + "search?key={keyword}&gorepo=false&gosubrepo=false&cmd=true&cgo=false"
	pkginfoApi = base + "pkginfo?pkgname={pkgname}"

	httpClient *http.Client
)

func init() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient = &http.Client{Transport: tr}
}

var (
	ErrPkgNotExists    = errors.New("gowalker: package not exist")
	ErrPkgNotGolangCmd = errors.New("gowalker: package not golang cmd package")
)

type PackageItem struct {
	Id          int64  `json:"id"`
	Name        string `json:"project_name"`
	Path        string `json:"project_path"`
	HomePage    string `json:"homepage"`
	ImportPath  string `json:"import_path"`
	IsCgo       bool   `json:"cgo"`
	IsCmd       bool   `json:"cmd"`
	Tags        string `json:"tags"` // like "master|||1.0.0RC1|||beego1.0.0|||", split by |||
	Description string `json:"synopsis"`
}

type SearchPackages struct {
	Packages []*PackageItem `json:"packages"`
}

func NewSearch(key string) (*SearchPackages, error) {
	url := com.Expand(searchApi, map[string]string{
		"keyword": key,
	})
	packages := new(SearchPackages)
	err := com.HttpGetJSON(httpClient, url, packages)
	return packages, err
}

func RefreshPkg(pkgname string) error {
	resp, err := http.Get("https://gowalker.org/" + pkgname)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func GetCmdPkgInfo(pkgname string) (*PackageItem, error) {
	pkginfo, err := GetPkgInfo(pkgname)
	if err != nil {
		return nil, err
	}
	if pkginfo.IsCmd == false {
		return nil, ErrPkgNotGolangCmd
	}
	return pkginfo, err
}

// https://gowalker.org/api/v1/pkginfo?pkgname=github.com/astaxie/beego
func GetPkgInfo(pkgname string) (*PackageItem, error) {
	err := RefreshPkg(pkgname)
	if err != nil {
		return nil, err
	}
	url := com.Expand(pkginfoApi, map[string]string{
		"pkgname": pkgname,
	})
	pkginfo := new(PackageItem)
	if err = com.HttpGetJSON(httpClient, url, pkginfo); err != nil {
		return nil, err
	}
	if pkginfo.Id == 0 {
		return nil, ErrPkgNotExists
	}
	return pkginfo, err
}
