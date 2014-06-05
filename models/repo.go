package models

import (
	"errors"
	"time"
)

var (
	ErrRepositoryNotExists = errors.New("repo not found")
)

type Repository struct {
	Id      int64
	Uri     string `xorm:"unique(r)"`
	Brief   string
	IsCgo   bool
	Created time.Time `xorm:"created"`
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
