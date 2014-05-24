package base

import (
	"errors"
	"strings"
)

func SanitizedRepoPath(repo string) string {
	repo = strings.TrimSpace(repo)
	if strings.HasSuffix(repo, ".git") {
		repo = repo[:len(repo)-4]
	}
	if strings.HasPrefix(repo, "https://") {
		repo = repo[len("https://"):]
	}
	return repo
}

type CVSInfo struct {
	Provide        string
	VersionControl string
	Owner          string
	Branch         string
	RepoName       string
	RepoSubPath    string
}

var (
	ErrCvsURIInvalid    = errors.New("cvs uri invalid")
	ErrCvsNotRecognized = errors.New("cvs path not recognized")
)

func ParseCvsURI(uri string) (*CVSInfo, error) {
	uri = SanitizedRepoPath(uri)
	fields := strings.Split(uri, "/")
	if len(fields) < 3 {
		return nil, ErrCvsURIInvalid
	}
	if strings.HasPrefix(uri, "github.com") {
		return &CVSInfo{
			Provide:        "github.com",
			VersionControl: "git",
			Owner:          fields[1],
			RepoName:       fields[2],
			RepoSubPath:    strings.Join(fields[2:], "/"),
		}, nil
	}
	return nil, ErrCvsNotRecognized
}
