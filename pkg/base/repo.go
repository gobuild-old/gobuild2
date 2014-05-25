package base

import (
	"errors"
	"strings"
)

func sanitizedRepoPath(repo string) string {
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
	Provider       string
	VersionControl string
	Owner          string
	Branch         string
	RepoName       string
	RepoSubPath    string
	FullPath       string
}

var (
	ErrCvsURIInvalid    = errors.New("cvs uri invalid")
	ErrCvsNotRecognized = errors.New("cvs path not recognized")
)

const (
	ProviderGithub = "github.com"
	ProviderGoogle = "code.google.com"
)

var verCtrlMap = map[string]string{
	ProviderGithub: "git",
	ProviderGoogle: "hg",
}

func ParseCvsURI(uri string) (*CVSInfo, error) {
	uri = sanitizedRepoPath(uri)
	var provider string
	switch {
	case strings.HasPrefix(uri, ProviderGoogle):
		provider = ProviderGoogle
	case strings.HasPrefix(uri, ProviderGithub):
		provider = ProviderGithub
	default:
		provider = ProviderGithub
		uri = ProviderGithub + "/" + uri
	}
	fields := strings.Split(uri, "/")
	if len(fields) < 3 {
		return nil, ErrCvsURIInvalid
	}
	return &CVSInfo{
		Provider:       provider,
		VersionControl: verCtrlMap[provider],
		Owner:          fields[1],
		RepoName:       fields[2],
		RepoSubPath:    strings.Join(fields[2:], "/"),
		FullPath:       uri,
	}, nil
	// cvsinfo := new(CVSInfo)
	// if strings.HasPrefix(uri, ProviderGithub) {
	// }
	// return nil, ErrCvsNotRecognized
}
