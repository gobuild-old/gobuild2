package models

import (
	"errors"
	"strings"
	"time"

	"github.com/gobuild/gobuild2/pkg/base"
	"github.com/gobuild/gobuild2/pkg/gowalker"
	"github.com/gobuild/log"
)

var (
	ErrRepositoryNotExists = errors.New("repo not found")
)

type Repository struct {
	Id            int64
	Uri           string `xorm:"unique(r)"`
	Brief         string
	IsCgo         bool
	IsCmd         bool
	Tags          []string
	Created       time.Time `xorm:"created"`
	DownloadCount int64
}

type RepoStatistic struct {
	Rid           int64 `xorm:"pk"`
	Pv            int64
	DownloadCount int64
	Updated       time.Time `xorm:"updated"`
}

type LastRepoUpdate struct {
	Rid        int64  `xorm:"unique(u)"`
	TagBranch  string `xorm:"unique(u)"`
	Os         string `xorm:"unique(u)"`
	Arch       string `xorm:"unique(u)"`
	PushURI    string
	ZipBallUrl string
	Updated    time.Time `xorm:"updated"`
}

func AddRepository(repoName string) (r *Repository, err error) {
	cvsinfo, err := base.ParseCvsURI(repoName) // base.SanitizedRepoPath(rf.Name)
	if err != nil {
		log.Errorf("parse cvs url error: %v", err)
		return
	}

	repoUri := cvsinfo.FullPath
	r = new(Repository)
	r.Uri = repoUri

	pkginfo, err := gowalker.GetPkgInfo(repoUri)
	if err != nil {
		log.Errorf("gowalker not passed check: %v", err)
		return
	}
	r.IsCgo = pkginfo.IsCgo
	r.IsCmd = pkginfo.IsCmd
	r.Tags = strings.Split(pkginfo.Tags, "|||")
	// description
	r.Brief = pkginfo.Description
	base.ParseCvsURI(repoUri)
	if strings.HasPrefix(repoUri, "github.com") {
		// comunicate with github
		fields := strings.Split(repoUri, "/")
		owner, repoName := fields[1], fields[2]
		repo, _, err := GHClient.Repositories.Get(owner, repoName)
		if err != nil {
			log.Errorf("get information from github error: %v", err)
		} else {
			r.Brief = *repo.Description
		}
	}
	if _, err = CreateRepository(r); err != nil {
		log.Errorf("create repo error: %v", err)
		return
	}
	return r, nil
}

func UpdateRepository(v *Repository, condi *Repository) (int64, error) {
	return orm.UseBool().Update(v, condi)
}

func CreateRepository(r *Repository) (*Repository, error) {
	// r := &Repository{Uri: repoUri}
	if has, err := orm.Get(r); err == nil && has {
		return r, nil
	}
	_, err := orm.Insert(r)
	return r, err
}

func GetAllRepos(count, start int) ([]Repository, error) {
	var rs []Repository
	err := orm.Limit(count, start).Desc("created").Find(&rs)
	return rs, err
}

func GetRepositoryById(id int64) (*Repository, error) {
	r := new(Repository)
	if has, err := orm.Id(id).Get(r); err == nil && has {
		return r, nil
	}
	return nil, ErrRepositoryNotExists
}

func GetRepositoryByName(name string) (*Repository, error) {
	r := &Repository{Uri: name}
	if has, err := orm.Get(r); err == nil && has {
		return r, nil
	}
	return nil, ErrRepositoryNotExists
}

func GetAllLastRepoByOsArch(os, arch string) (us []LastRepoUpdate, err error) {
	err = orm.Asc("rid").Find(&us, &LastRepoUpdate{Os: os, Arch: arch})
	return us, err
}

func GetAllLastRepoUpdate(rid int64) (us []LastRepoUpdate, err error) {
	err = orm.Find(&us, &LastRepoUpdate{Rid: rid})
	return
}
