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
	if strings.HasPrefix(repo, "http://") {
		repo = repo[len("http://"):]
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
	PVD_GITHUB      = "github.com"
	PVD_GOOGLE      = "code.google.com"
	defaultProvider = PVD_GITHUB
)

var ProviderCtrlMap = map[string]string{
	PVD_GITHUB: "git",
	PVD_GOOGLE: "hg",
}

func ParseCvsURI(uri string) (*CVSInfo, error) {
	uri = sanitizedRepoPath(uri)
	var provider string
	guessProvider := strings.Split(uri, "/")[0]
	if _, has := ProviderCtrlMap[guessProvider]; has {
		provider = guessProvider
	}
	if provider == "" {
		provider = defaultProvider
		uri = provider + "/" + uri
	}

	fields := strings.Split(uri, "/")
	if len(fields) < 3 {
		return nil, ErrCvsURIInvalid
	}
	branch := ""
	switch provider {
	case PVD_GITHUB:
		branch = "master"
	case PVD_GOOGLE:
		branch = "default" // for hg
	}
	// log.Infof("branch: %v", branch)
	return &CVSInfo{
		Provider:       provider,
		Branch:         branch,
		VersionControl: ProviderCtrlMap[provider],
		Owner:          fields[1],
		RepoName:       fields[2],
		RepoSubPath:    strings.Join(fields[2:], "/"),
		FullPath:       uri,
	}, nil
}
