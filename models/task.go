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
	Uri   string
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
	Branch    string // can also be tag, commit_id
	Os        string
	Arch      string
	CgoEnable bool

	Status  string
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}

var ErrTaskNotAvaliable = errors.New("task not avaliable for now")

func init() {
	tables = append(tables, new(Task), new(Repository), new(RepoStatistic), new(DownloadHistory))
}

func CreateTask(task *Task) (*Task, error) {
	task.Status = ST_READY
	_, err := orm.Insert(task)
	return task, err
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
