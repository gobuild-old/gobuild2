package github

import (
	"time"

	"github.com/google/go-github/github"
)

// https://developer.github.com/v3/repos/#get
type RepoItem struct {
	Name        string    `json:"name"`
	FullName    string    `json:"full_name"`
	Description string    `json:"description"`
	Private     bool      `json:"private"`
	Updated     time.Time `json:"updated_at"`
}

const api = "https://api.github.com/"

func GetRepoInfo(owner string, repoName string) (*github.Repository, error) {
	// base.HttpGetJSON(api + "/repos")
	client := github.NewClient(nil)
	repo, _, err := client.Repositories.Get(owner, repoName)
	return repo, err
}
