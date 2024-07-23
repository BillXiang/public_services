package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"public_service/config"
	public_service "public_service/service"
	"syscall"

	"github.com/go-git/go-git/v5"
	"github.com/jacobsa/daemonize"
)

var (
	configForeground = flag.Bool("f", false, "run foreground")
	configService    = flag.Bool("s", false, "start self as service")
	configFile       = flag.String("c", "./public_service.conf", "config file")
)

func main() {
	flag.Parse()
	cfg, cfgErr := NewClientCfg(*configFile)
	if cfgErr != nil {
		fmt.Printf("Critical error happened %s, try again", cfgErr.Error())
		_ = daemonize.SignalOutcome(cfgErr)
		os.Exit(0)
	}
	git_url := cfg.GetGit()
	path := cfg.GetGitLocal()

	if !*configForeground {
		log.Println("!configForeground")
		if err := startDaemon(path); err != nil {
			fmt.Printf("startDaemon failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Start sueecss, exit\n")
		os.Exit(0)
	}

	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:               git_url,
		RecurseSubmodules: git.NoRecurseSubmodules,
	})
	if err == nil || err.Error() == "repository already exists" {
		if *configService {
			log.Println("StartSelfAsService")
			public_service.StartSelfAsService(path)
		} else {
			r, err := git.PlainOpen(path)
			if err != nil {
				fmt.Println(err.Error())
				_ = daemonize.SignalOutcome(err)
				return
			}
			w, err := r.Worktree()
			if err != nil {
				fmt.Println("Worktree " + err.Error())
				_ = daemonize.SignalOutcome(err)
				return
			}
			err = w.Pull(&git.PullOptions{RemoteName: "origin"})
			if err == nil {
				fmt.Println("Pull success")
			} else {
				fmt.Println("Pull " + err.Error())
			}

			jobCfg, err := config.LoadConfigFile(path + "/public_service.json")
			if err != nil {
				_ = daemonize.SignalOutcome(err)
				return
			}

			public_service.Exec(jobCfg)
			public_service.Cron(jobCfg)
			public_service.Service(jobCfg, path)
		}
		_ = daemonize.SignalOutcome(nil)

		sigs := make(chan os.Signal, 1)
		done := make(chan bool, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			sig := <-sigs
			fmt.Println()
			fmt.Println(sig)
			done <- true
		}()
		fmt.Println("awaiting signal")
		<-done
		fmt.Println("exiting")
	}
	_ = daemonize.SignalOutcome(err)
}

func startDaemon(path string) error {
	cmdPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("startDaemon failed: cannot get absolute command path, err(%v)", err)
	}

	args := []string{"-s", "-f"}
	args = append(args, os.Args[1:]...)

	env := os.Environ()

	// add GODEBUG=madvdontneed=1 environ, to make sysUnused uses madvise(MADV_DONTNEED) to signal the kernel that a
	// range of allocated memory contains unneeded data.
	env = append(env, "GODEBUG=madvdontneed=1")
	env = append(env, "PMON2_CONF="+path+"/pmon2_conf.yml")
	err = daemonize.Run(cmdPath, args, env, os.Stdout)
	if err != nil {
		return fmt.Errorf("startDaemon failed: daemon start failed, cmd(%v) args(%v) env(%v) err(%v)", cmdPath, args, env, err)
	}
	return nil
}
