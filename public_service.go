package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"public_service/daemon"
	public_service "public_service/service"

	"github.com/cubefs/cubefs/util/config"
	"github.com/go-git/go-git/v5"
	"github.com/jacobsa/daemonize"
	service "github.com/ntt360/pmon2/app"
	service_monitor "github.com/ntt360/pmon2/app/god"
	"github.com/ntt360/pmon2/app/model"
	"github.com/ntt360/pmon2/app/output"
	"github.com/robfig/cron/v3"
)

var (
	configForeground = flag.Bool("f", false, "run foreground")
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
	log.Println("configForeground")
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:               git_url,
		RecurseSubmodules: git.NoRecurseSubmodules,
	})
	if err == nil || err.Error() == "repository already exists" {
		r, err := git.PlainOpen(path)
		if err != nil {
			fmt.Println(err.Error())
			_ = daemonize.SignalOutcome(err)
			return
		}
		w, err := r.Worktree()
		if err != nil {
			fmt.Println(err.Error())
			_ = daemonize.SignalOutcome(err)
			return
		}
		err = w.Pull(&git.PullOptions{RemoteName: "origin"})
		if err == nil {
			fmt.Println("Pull success")
		} else {
			fmt.Println(err.Error())
		}

		jobCfg, err := config.LoadConfigFile(path + "/public_service.json")
		if err != nil {
			_ = daemonize.SignalOutcome(err)
			return
		}

		execJobs := jobCfg.GetSlice("exec")
		for _, item := range execJobs {
			execJob := item.(map[string]interface{})
			paras := execJob["para"].([]interface{})
			parasStrs := make([]string, 0, len(paras))
			for _, para := range paras {
				parasStrs = append(parasStrs, para.(string))
			}
			pd := &daemon.ProcessDaemon{
				CmdPath: execJob["cmd"].(string),
				Args:    parasStrs,
			}
			pd.StartSubprocess()
		}

		cron := cron.New(cron.WithSeconds())
		cronJobs := jobCfg.GetSlice("cron")
		for _, item := range cronJobs {
			cronJob := item.(map[string]interface{})
			paras := cronJob["para"].([]interface{})
			parasStrs := make([]string, 0, len(paras))
			for _, para := range paras {
				parasStrs = append(parasStrs, para.(string))
			}
			pd := &daemon.ProcessDaemon{
				CmdPath: cronJob["cmd"].(string),
				Args:    parasStrs,
			}
			cron.AddFunc(cronJob["exp"].(string), func() {
				pd.StartSubprocess()
			})
		}
		cron.Start()

		err = service.Instance(path + "/pmon2_conf.yml")
		if err != nil {
			_ = daemonize.SignalOutcome(err)
			log.Fatal(err)
		}
		serviceJobs := jobCfg.GetSlice("service")
		for _, item := range serviceJobs {
			var flag model.ExecFlags
			serviceJob := item.(map[string]interface{})
			if serviceJob["no_auto_restart"] != nil {
				flag.NoAutoRestart = (serviceJob["no_auto_restart"].(string) == "true")
			}
			if serviceJob["args"] != nil {
				flag.Args = serviceJob["args"].(string)
			}
			serviceRun([]string{serviceJob["cmd"].(string)}, flag, path)
		}
		_ = daemonize.SignalOutcome(nil)
		// start monitor service
		service_monitor.NewMonitor()
	}
	_ = daemonize.SignalOutcome(err)
}

func serviceRun(args []string, flag model.ExecFlags, path string) {
	// get exec abs file path
	execPath, err := public_service.GetExecFile(args)
	if err != nil {
		service.Log.Error(err.Error())
		return
	}
	flag.Log = path + "/pmon2/log/" + filepath.Base(execPath) + ".log"
	flags := flag.Json()

	m, exist := public_service.ProcessExist(execPath)
	var rel []string
	if exist {
		service.Log.Debugf("restart process: %v", flags)
		rel, err = public_service.Restart(m, flags)
	} else {
		service.Log.Debugf("load first process: %v", flags)
		rel, err = public_service.LoadFirst(execPath, flags)
	}

	if err != nil {
		if len(os.Getenv("PMON2_DEBUG")) > 0 {
			service.Log.Debugf("%+v", err)
		} else {
			service.Log.Debugf(err.Error())
		}
	}
	if rel != nil {
		output.TableOne(rel)
	}
}

func startDaemon(path string) error {
	cmdPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("startDaemon failed: cannot get absolute command path, err(%v)", err)
	}

	args := []string{"-f"}
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
