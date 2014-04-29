package models

import "time"

const (
	TS_PENDING = "pending"
	TS_START   = "start"
	TS_DONE    = "done"
)

type Task struct {
	Id       int64
	RepoAddr string
	Status   string
	Created  time.Time `xorm:"created"`
	Updated  time.Time `xorm:"updated"`
}

func init() { tables = append(tables, new(Task)) }

func CreateTask(task *Task) (*Task, error) {
	_, err := x.Insert(task)
	return task, err
}

func UpdateTaskStatus(tid int64, status string) error {
	_, err := x.Id(tid).Update(&Task{Status: status})
	return err
}
