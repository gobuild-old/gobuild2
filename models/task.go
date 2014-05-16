package models

import (
	"errors"
	"time"
)

const (
	ST_READY      = "ready"
	ST_PENDING    = "pending"
	ST_RETRIVING  = "retriving"
	ST_BUILDING   = "building"
	ST_PUBLISHING = "publishing"
	ST_DONE       = "done"
	ST_ERROR      = "error"
)

type Repository struct {
	Id    int64
	Uri   string `xorm:"unique(r)"`
	Brief string
}

type RepoStatistic struct {
	Rid           int64 `xorm:"pk"`
	Pv            int64
	DownloadCount int64
	Updated       time.Time `xorm:"updated"`
}

type DownloadHistory struct {
	Id       int64
	Rid      int64
	Os       string
	Arch     string
	Ip       string
	CurrTime time.Time `xorm:"created"`
}

type Task struct {
	Id        int64
	Rid       int64
	Repo      *Repository `xorm:"-"`
	Branch    string      // can also be tag, commit_id
	Os        string
	Arch      string
	CgoEnable bool

	Status  string
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}

var (
	ErrTaskNotAvaliable    = errors.New("not task ready for build now")
	ErrTaskNotExists       = errors.New("task not exists")
	ErrRepositoryNotExists = errors.New("repo not found")
)

func init() {
	tables = append(tables, new(Task), new(Repository), new(RepoStatistic), new(DownloadHistory))
}

func CreateRepository(repoUri string) (*Repository, error) {
	r := &Repository{Uri: repoUri}
	if has, err := orm.Get(r); err == nil && has {
		return r, nil
	}
	r.Uri = repoUri
	r.Brief = "todo, not get"
	_, err := orm.Insert(r)
	return r, err
}

func GetRepositoryById(id int64) (*Repository, error) {
	r := new(Repository)
	if has, err := orm.Id(id).Get(r); err == nil && has {
		return r, nil
	}
	return nil, ErrRepositoryNotExists
}

func CreateTask(task *Task) (*Task, error) {
	task.Status = ST_READY
	_, err := orm.Insert(task)
	return task, err
}

func GetTaskById(tid int64) (*Task, error) {
	t := new(Task)
	has, err := orm.Id(tid).Get(t)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrTaskNotExists
	}
	return t, nil
}

func ResetAllTaskStatus() error {
	_, err := orm.Where("status != ?", ST_DONE).Update(&Task{Status: ST_READY})
	return err
}

func UpdateTaskStatus(tid int64, status string) error {
	_, err := orm.Id(tid).Update(&Task{Status: status})
	return err
}

func GetAvaliableTask(cgo bool, os, arch string) (*Task, error) {
	task := &Task{Status: ST_READY}
	if cgo == false {
		task.Os = os
		task.Arch = arch
	}
	exists, err := orm.Get(task)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrTaskNotAvaliable
	}
	return task, nil
}
