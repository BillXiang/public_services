package public_service

import (
	"fmt"
	"public_service/config"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func RemoteProcess(repo *git.Repository, path string) {
	go func() {
		var since *time.Time = nil
		w, err := repo.Worktree()
		if err != nil {
			fmt.Println("Worktree " + err.Error())
			return
		}
		for {
			fmt.Println(since)
			err := w.Pull(&git.PullOptions{})
			if err == nil {
				fmt.Println("Pull success")
			} else {
				fmt.Println("Pull " + err.Error())
			}

			commits, _ := repo.Log(&git.LogOptions{
				Since: since,
			})
			commits.ForEach(func(c *object.Commit) error {
				if since == nil || c.Author.When.After(*since) {
					since = &c.Author.When
					fmt.Print(c)
					if strings.HasPrefix(c.Message, "[REMOTE_PROCESS]") {
						jobCfg, err := config.LoadConfigFile(path + "/remote_process.json")
						if err != nil {
							return err
						}
						Exec(jobCfg)
					}
				}
				return nil
			})

			time.Sleep(10 * time.Second)
		}
	}()
}
