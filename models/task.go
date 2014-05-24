package models

import (
	"errors"
	"strings"
	"time"

	"code.google.com/p/goauth2/oauth"
	"github.com/gobuild/gobuild2/pkg/base"
	"github.com/gobuild/gobuild2/pkg/gowalker"
	"github.com/google/go-github/github"
	"github.com/qiniu/log"
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

type DownloadHistory struct {
	Id       int64
	Rid      int64
	Os       string
	Arch     string
	Ip       string
	CurrTime time.Time `xorm:"created"`
}

type BuildHistory struct {
	Id      int64
	Tid     int64  `xorm:"unique(b)"`
	Status  string `xorm:"unique(b)"`
	Output  string `xorm:"TEXT"`
	Updated string `xorm:"updated"`
}

type Task struct {
	Id            int64
	Rid           int64       `xorm:"unique(t)"`
	Repo          *Repository `xorm:"-"`
	Os            string      `xorm:"unique(t)"`
	Arch          string      `xorm:"unique(t)"`
	CgoEnable     bool
	CommitMessage string
	ArchieveAddr  string

	Branch string
	Sha    string `xorm:"unique(t)"`

	Status  string
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}

type LastRepoUpdate struct {
	Rid        int64  `xorm:"unique(u)"`
	Branch     string `xorm:"unique(u)"`
	OsArch     string `xorm:"unique(u)"`
	Sha        string
	ZipBallUrl string
	Updated    time.Time `xorm:"updated"`
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

var gclient *github.Client //= github.NewClient(httpClient)
func init() {
	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: githubPublicAccessToken},
	}
	gclient = github.NewClient(t.Client())
}

func CreateRepository(repoUri string) (*Repository, error) {
	pkginfo, err := gowalker.GetCmdPkgInfo(repoUri)
	if err != nil {
		log.Errorf("gowalker not passed check: %v", err)
		return nil, err
	}
	r := &Repository{Uri: repoUri}
	if has, err := orm.Get(r); err == nil && has {
		return r, nil
	}
	r.Uri = repoUri
	r.IsCgo = pkginfo.IsCgo
	// description
	r.Brief = pkginfo.Description
	if strings.HasPrefix(repoUri, "github.com") {
		// comunicate with github
		fields := strings.Split(repoUri, "/")
		owner, repoName := fields[1], fields[2]
		repo, _, err := gclient.Repositories.Get(owner, repoName)
		if err != nil {
			log.Errorf("get information from github error: %v", err)
		} else {
			r.Brief = *repo.Description
		}
	}
	_, err = orm.Insert(r)
	return r, err
}

func CreateNewBuilding(rid int64, branch string, os, arch string) (err error) {
	repo, err := GetRepositoryById(rid)
	if err != nil {
		return
	}
	task := &Task{
		Rid:       rid,
		Branch:    branch,
		Os:        os,
		Arch:      arch,
		CgoEnable: repo.IsCgo,
	}
	var cvsinfo *base.CVSInfo
	if cvsinfo, err = base.ParseCvsURI(repo.Uri); err != nil {
		return
	}
	info, _, err := gclient.Repositories.GetBranch(cvsinfo.Owner, cvsinfo.RepoName, branch)
	if err != nil {
		return
	}
	task.Sha = *info.Commit.SHA
	if info.Commit.Message != nil {
		task.CommitMessage = *info.Commit.Message // info.Commit.String() //"" // *info.Commit.Message
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

func CreateTask(task *Task) (*Task, error) {
	task.Status = ST_READY
	_, err := orm.Insert(task)
	return task, err
}

func GetTasksByRid(rid int64) ([]Task, error) {
	var ts []Task
	err := orm.Find(&ts, &Task{Rid: rid})
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

func UpdateTaskStatus(tid int64, status string, output string) error {
	pubAddr := ""
	if status == ST_DONE {
		pubAddr = output
		tk, _ := GetTaskById(tid)
		condi := LastRepoUpdate{
			Rid:    tk.Rid,
			Branch: tk.Branch,
			OsArch: tk.Os + "-" + tk.Arch}
		lr := condi
		if has, err := orm.Get(&lr); err == nil && has {
			orm.Update(&LastRepoUpdate{Sha: tk.Sha, ZipBallUrl: pubAddr}, &condi)
		} else {
			condi.ZipBallUrl = pubAddr
			condi.Sha = tk.Sha
			if _, err := orm.Insert(&condi); err != nil {
				log.Errorf("insert last_repo_update failed: %v", err)
			}
		}
	}
	if _, err := orm.Id(tid).Update(&Task{Status: status, ArchieveAddr: pubAddr}); err != nil {
		return err
	}
	condi := &BuildHistory{Tid: tid, Status: status}
	if has, err := orm.Get(condi); err == nil && has {
		_, er := orm.Update(&BuildHistory{Output: output}, condi)
		return er
	}
	_, err := orm.Insert(&BuildHistory{Tid: tid, Status: status, Output: output})
	return err
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
		if exists, err := orm.Asc("created").Get(task); err != nil || !exists {
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
