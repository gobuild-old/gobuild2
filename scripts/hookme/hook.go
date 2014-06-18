/*{
	"name": "gobuild2",
	"secret": "123456",
	"repo": "github.com/gobuild/gobuild2",
	"zipball_url": "http://xxx.zip"
}
*/

package main

type GogsUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type GogsCommit struct {
	Id      string   `json:"id"`
	Message string   `json:"message"`
	Url     string   `json:"url"`
	Author  GogsUser `json:"author"`
}

type GogsPayload struct {
	Secret  string       `json:"secret"`
	Ref     string       `json:"ref"`
	Commits []GogsCommit `json:"commits"`
	Pusher  GogsUser     `json:"pusher"`
}
