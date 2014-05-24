package base

import "strings"

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
