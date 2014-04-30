package social

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"code.google.com/p/goauth2/oauth"

	"github.com/gobuild/gobuild2/models"
	"github.com/gobuild/gobuild2/pkg/config"
)

type BasicUserInfo struct {
	Identity string
	Name     string
	Email    string
}

type Connector interface {
	Type() int
	SetRedirectUrl(string)
	UserInfo(*oauth.Token, *url.URL) (*BasicUserInfo, error)

	AuthCodeURL(string) string
	Exchange(string) (*oauth.Token, error)
}

var (
	SocialBaseUrl = "/login"
	SocialMap     = make(map[string]Connector)
)

func NewOauthService() {
	cfg := config.Config.Social
	for _, name := range []string{"github"} {
		if c, ok := cfg[name]; ok {
			oacf := &oauth.Config{
				ClientId:     c.ClientId,
				ClientSecret: c.ClientSecret,
			}
			oacf.AuthURL = c.AuthURL
			oacf.TokenURL = c.TokenURL
			oacf.Scope = SocialBaseUrl + "/" + name
			// oacf.RedirectURL = "todo"

			SocialMap[name] = &SocialGithub{
				Transport: &oauth.Transport{
					Config:    oacf,
					Transport: http.DefaultTransport,
				},
			}
		}
	}
}

//   ________.__  __                   __.
//  /  _____/|__|/  |_  /   |   \ __ __| |__
// /   \  ___|  \   __\ | _ ~ _ | | |  | __ \
// \    \_\  \  ||  |   |   Y   | | |  | \_\ \
//  \______  /__||__|   \   |   / |___/|___  /
//         \/                              \/

type SocialGithub struct {
	Token *oauth.Token
	*oauth.Transport
}

// func enableGitHub(config *oauth.Config) {
// 	SocialMap["github"] = &SocialGithub{
// 		Transport: &oauth.Transport{
// 			Config:    config,
// 			Transport: http.DefaultTransport,
// 		},
// 	}
// }

func (s *SocialGithub) Type() int                 { return models.OT_GITHUB }
func (s *SocialGithub) SetRedirectUrl(url string) { s.Transport.Config.RedirectURL = url }

func (s *SocialGithub) UserInfo(token *oauth.Token, _ *url.URL) (*BasicUserInfo, error) {
	transport := &oauth.Transport{
		Token: token,
	}
	var data struct {
		Id    int    `json:"id"`
		Name  string `json:"login"`
		Email string `json:"email"`
	}
	var err error
	r, err := transport.Client().Get(s.Transport.Scope)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	if err = json.NewDecoder(r.Body).Decode(&data); err != nil {
		return nil, err
	}
	return &BasicUserInfo{
		Identity: strconv.Itoa(data.Id),
		Name:     data.Name,
		Email:    data.Email,
	}, nil
}
