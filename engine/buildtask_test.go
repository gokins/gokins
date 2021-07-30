package engine

import (
	"fmt"
	"github.com/gokins/core/runtime"
	"testing"
)

func TestGitClone(t *testing.T) {
	task := BuildTask{
		build: &runtime.Build{
			Id: "1231",
			Repo: &runtime.Repository{
				Name:     "",
				Token:    "",
				Sha:      "c202ee042db1fc8b8c16c6c968195cec6185d7db",
				CloneURL: "https://gitee.com/SuperHeroJim/gokins-test",
			},
		},
	}
	err := task.getRepo()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(task.repoPath)
	fmt.Println(task.isClone)
}
