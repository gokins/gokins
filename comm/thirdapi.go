package comm

import (
	"github.com/gokins-main/gokins/thirdapi"
	"github.com/gokins-main/gokins/thirdapi/giteaapi"
	"github.com/gokins-main/gokins/thirdapi/giteeapi"
	"github.com/gokins-main/gokins/thirdapi/giteepremiumapi"
	"github.com/gokins-main/gokins/thirdapi/githubapi"
	"github.com/gokins-main/gokins/thirdapi/gitlabapi"
	"github.com/sirupsen/logrus"
)

var (
	apiClient *thirdapi.Client
)

func GetThirdApi(s string, host string) (*thirdapi.Client, error) {
	if apiClient == nil {
		switch s {
		case "gitee":
			apiClient = giteeapi.NewDefault()
		case "github":
			apiClient = githubapi.NewDefault()
		case "gitalb":
			client, err := gitlabapi.New(host + "/api/v4")
			if err != nil {
				return nil, err
			}
			apiClient = client
		case "giteepremium":
			client, err := giteepremiumapi.New(host + "/api/v5")
			if err != nil {
				return nil, err
			}
			apiClient = client
		case "gitea":
			client, err := giteaapi.New(host + "/api/v1")
			if err != nil {
				return nil, err
			}
			apiClient = client
		default:
			logrus.Debug("GetThirdApi default : 'github' ")
			apiClient = githubapi.NewDefault()
		}
	}
	return apiClient, nil
}
