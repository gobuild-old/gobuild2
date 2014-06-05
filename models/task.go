package models

import (
	"errors"
	"time"

	"code.google.com/p/goauth2/oauth"
	"github.com/gobuild/gobuild2/pkg/base"
	"github.com/gobuild/log"
	"github.com/google/go-github/github"
)

const (
	ST_READY      = "ready"
	ST_PENDING    = "pending"
	ST_RETRIVING  = "retriving"
	ST_BUILDING   = "building"
	ST_PACKING    = "packing"
	ST_PUBLISHING = "publishing"
	ST_DONE       = "done"
	ST_ERROR      = "error"
)

const (
	AC_BUILD  = "action-build"
	AC_SRCPKG = "action-source-package"
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

type CompileHistory struct {
	Id        int64
	CompileId int64
	Status    string
	Output    string `xorm:"TEXT"`
	Updated   string `xorm:"updated"`
}

type DownloadHistory struct {
	Id       int64
	Rid      int64
	Os       string
	Arch     string
	Ip       string
	CurrTime time.Time `xorm:"created"`
}

// going to clean
type BuildHistory struct {
	Id      int64
	Tid     int64  `xorm:"unique(b)"`
	Status  string `xorm:"unique(b)"`
	Output  string `xorm:"TEXT"`
	Updated string `xorm:"updated"`
}

// going to clean
type Task struct {
	Id            int64
	Rid           int64       `xorm:"unique(t)"`
	Repo          *Repository `xorm:"-"`
	Action        string      `xorm:"unique(t)"` // build or package
	Os            string      `xorm:"unique(t)"`
	Arch          string      `xorm:"unique(t)"`
	CgoEnable     bool
	CommitMessage string
	ZipBallUrl    string

	TagBranch string
	PushType  string //`xorm:"unique(t)"` // branch|tag|commit
	PushValue string `xorm:"unique(t)"` //master|v1.2|45asf913

	Status  string
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
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

func GetAllLastRepoByOsArch(os, arch string) (us []LastRepoUpdate, err error) {
	err = orm.Asc("rid").Find(&us, &LastRepoUpdate{Os: os, Arch: arch})
	return us, err
}

func GetAllLastRepoUpdate(rid int64) (us []LastRepoUpdate, err error) {
	err = orm.Find(&us, &LastRepoUpdate{Rid: rid})
	return
}

var (
	ErrTaskNotAvaliable    = errors.New("not task ready for build now")
	ErrTaskNotExists       = errors.New("task not exists")
	ErrRepositoryNotExists = errors.New("repo not found")
	ErrNoAvaliableDownload = errors.New("no avaliable download")
	ErrTaskIsRunning       = errors.New("task is running")
)

func init() {
	tables = append(tables, new(Task),
		new(Repository), new(RepoStatistic), new(LastRepoUpdate),
		new(DownloadHistory), new(BuildHistory))
}

const githubPublicAccessToken = "68655dcb3723d24e2fbe2cb450747c14b966eac3"

var GHClient *github.Client //= github.NewClient(httpClient)
func init() {
	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: githubPublicAccessToken},
	}
	GHClient = github.NewClient(t.Client())
}

func CreateRepository(r *Repository) (*Repository, error) {
	// r := &Repository{Uri: repoUri}
	if has, err := orm.Get(r); err == nil && has {
		return r, nil
	}
	_, err := orm.Insert(r)
	return r, err
}

func CreateNewBuilding(rid int64, branch string, os, arch string, action string) (err error) {
	repo, err := GetRepositoryById(rid)
	if err != nil {
		return
	}
	task := &Task{
		Rid:       rid,
		Action:    action,
		TagBranch: branch,
		PushType:  "branch",
		PushValue: branch,
		Os:        os,
		Arch:      arch,
		CgoEnable: repo.IsCgo,
	}
	var cvsinfo *base.CVSInfo
	log.Infof("add task for repo: %v", repo.Uri)
	if cvsinfo, err = base.ParseCvsURI(repo.Uri); err != nil {
		return
	}
	if cvsinfo.Provider == base.PVD_GITHUB {
		info, _, er := GHClient.Repositories.GetBranch(cvsinfo.Owner, cvsinfo.RepoName, branch)
		if er != nil {
			err = er
			return
		}
		log.Infof("get information from github:%v", info)
		task.PushType = "commit"
		task.PushValue = *info.Commit.SHA
		if info.Commit.Message != nil {
			task.CommitMessage = *info.Commit.Message
		}
	}
	_, err = CreateTask(task)
	return
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

func CreateTasks(tasks []*Task) error {
	for _, t := range tasks {
		t.Status = ST_READY
	}
	_, err := orm.Insert(tasks)
	return err
}

func CreateTask(task *Task) (*Task, error) {
	task.Status = ST_READY
	_, err := orm.Insert(task)
	return task, err
}

func GetTasksByRid(rid int64) ([]Task, error) {
	var ts []Task
	err := orm.Desc("id").Find(&ts, &Task{Rid: rid})
	return ts, err
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
	t.Repo, err = GetRepositoryById(t.Rid)
	return t, err
}

func GetOneDownloadableTask(rid int64, os, arch string) (*Task, error) {
	task := &Task{Os: os, Arch: arch, Rid: rid, Status: ST_DONE}
	ok, err := orm.Desc("id").Get(task)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrNoAvaliableDownload
	}
	return task, nil
}

func ResetTask(tid int64) error {
	task, err := GetTaskById(tid)
	if err != nil {
		return err
	}
	if task.Status == ST_DONE || task.Status == ST_ERROR {
		orm.Delete(&BuildHistory{Tid: task.Id})
		orm.Id(task.Id).Update(&Task{Status: ST_READY})
	}
	return nil
}

func ResetAllTaskStatus() error {
	_, err := orm.Where("status != ? and status != ?", ST_DONE, ST_ERROR).Update(&Task{Status: ST_READY})
	return err
}

func UpdatePubAddr(tid int64, pubAddr string) error {
	tk, err := GetTaskById(tid)
	if err != nil {
		return err
	}
	condi := LastRepoUpdate{
		Rid:       tk.Rid,
		TagBranch: tk.TagBranch,
		Os:        tk.Os,
		Arch:      tk.Arch}
	lr := condi
	pushURI := tk.PushType + ":" + tk.PushValue
	if has, err := orm.Get(&lr); err == nil && has {
		orm.Update(&LastRepoUpdate{PushURI: pushURI, ZipBallUrl: pubAddr}, &condi)
	} else {
		condi.ZipBallUrl = pubAddr
		condi.PushURI = pushURI
		if _, err := orm.Insert(&condi); err != nil {
			log.Errorf("insert last_repo_update failed: %v", err)
		}
	}
	if _, err := orm.Id(tid).Update(&Task{ZipBallUrl: pubAddr}); err != nil {
		return err
	}
	return nil
}

func UpdateTaskStatus(tid int64, status string, output string) error {
	// pubAddr := ""
	// if status == ST_DONE {
	// 	pubAddr = output
	// 	tk, _ := GetTaskById(tid)
	// 	condi := LastRepoUpdate{
	// 		Rid:       tk.Rid,
	// 		TagBranch: tk.TagBranch,
	// 		Os:        tk.Os,
	// 		Arch:      tk.Arch}
	// 	lr := condi
	// 	pushURI := tk.PushType + ":" + tk.PushValue
	// 	if has, err := orm.Get(&lr); err == nil && has {
	// 		orm.Update(&LastRepoUpdate{PushURI: pushURI, ZipBallUrl: pubAddr}, &condi)
	// 	} else {
	// 		condi.ZipBallUrl = pubAddr
	// 		condi.PushURI = pushURI
	// 		if _, err := orm.Insert(&condi); err != nil {
	// 			log.Errorf("insert last_repo_update failed: %v", err)
	// 		}
	// 	}
	// }
	// if _, err := orm.Id(tid).Update(&Task{Status: status, ArchieveAddr: pubAddr}); err != nil {
	// 	return err
	// }
	log.Debugf("update task(%d) status(%s)", tid, status)
	if _, err := orm.Id(tid).Update(&Task{Status: status}); err != nil {
		log.Errorf("update task status error: %v", err)
	}
	condi := &BuildHistory{Tid: tid, Status: status}
	if has, err := orm.Get(condi); err == nil && has {
		_, er := orm.Update(&BuildHistory{Output: output}, condi)
		return er
	}
	_, err := orm.Insert(&BuildHistory{Tid: tid, Status: status, Output: output})
	return err
}

func GetAvaliableTasks(os, arch string) (tasks []*Task, err error) {
	t, e := GetAvaliableTask(os, arch)
	if e != nil {
		return nil, e
	}
	if t.CgoEnable {
		return []*Task{t}, nil
	}
	return []*Task{t}, nil
}

func GetAvaliableTask(os, arch string) (task *Task, err error) {
	task = &Task{Status: ST_READY, CgoEnable: false}
	defer func() {
		if task != nil && task.Id != 0 {
			task, err = GetTaskById(task.Id)
		}
	}()
	exists, err := orm.UseBool().Asc("created").Get(task)
	if err != nil {
		return nil, err
	}
	if !exists {
		task.CgoEnable = true
		task.Os, task.Arch = os, arch
		if exists, err := orm.UseBool().Asc("created").Get(task); err != nil || !exists {
			return nil, ErrTaskNotAvaliable
		}
	}
	affec, err := orm.Id(task.Id).Update(&Task{Status: ST_PENDING}, &Task{Status: ST_READY})
	if err != nil {
		return nil, err
	}
	if affec == 0 { // task already taken away by another process
		return nil, ErrTaskNotAvaliable
	}
	return task, nil
}

func GetAllBuildHistoryByTid(tid int64) ([]BuildHistory, error) {
	var bh []BuildHistory
	err := orm.Asc("id").Find(&bh, &BuildHistory{Tid: tid})
	return bh, err
}
