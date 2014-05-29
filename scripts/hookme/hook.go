/*{
	"name": "gobuild2",
	"secret": "123456",
	"repo": "github.com/gobuild/gobuild2",
	"zipball_url": "http://xxx.zip"
}
*/

package main

type HookInfo struct {
	Name       string `json:"name"`
	Secret     string `json:"secret"`
	Repo       string `json:"repo"`
	ZipballUrl string `json:"zipball_url"`
}
