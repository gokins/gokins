package engine

import (
	"context"
	"fmt"
	"github.com/gokins/core/runtime"
	"testing"
)

func TestGitClone(t *testing.T) {
	task := BuildTask{
		repoPath: "E:\\workspace\\tst\\abc",
		repoPaths: "E:\\workspace\\tst\\abc",
		isClone: true,
		ctx: context.Background(),
		build: &runtime.Build{
			Id: "1231",
			Repo: &runtime.Repository{
				Name:     "",
				Token:    "",
				Sha:      "testing",
				CloneURL: "https://foot/foo.git",
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
