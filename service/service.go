package public_service

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"public_service/config"
	"strings"

	"github.com/jacobsa/daemonize"
	service "github.com/ntt360/pmon2/app"
	service_monitor "github.com/ntt360/pmon2/app/god"
	"github.com/ntt360/pmon2/app/model"
	"github.com/ntt360/pmon2/app/output"
)

func Service(jobCfg *config.Config, path string) {
	err := service.Instance(path + "/pmon2_conf.yml")
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
		ServiceRun([]string{serviceJob["cmd"].(string)}, flag, path)
	}
	// start monitor service
	go func() {
		service_monitor.NewMonitor()
	}()

}

func ServiceRun(args []string, flag model.ExecFlags, path string) {
	// get exec abs file path
	execPath, err := GetExecFile(args)
	if err != nil {
		service.Log.Error(err.Error())
		return
	}
	dir := path + "/pmon2/log/"
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			service.Log.Errorf("err: %s, logs dir: '%s'", err.Error(), dir)
		}
	}
	flag.Log = dir + filepath.Base(execPath) + ".log"
	flags := flag.Json()

	m, exist := ProcessExist(execPath)
	var rel []string
	if exist {
		service.Log.Debugf("restart process: %v", flags)
		rel, err = Restart(m, flags)
	} else {
		service.Log.Debugf("load first process: %v", flags)
		rel, err = LoadFirst(execPath, flags)
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

func StartSelfAsService(path string) error {
	cmdPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("StartSelfAsService failed: cannot get absolute command path, err(%v)", err)
	}
	err = service.Instance(path + "/pmon2_self_conf.yml")
	if err != nil {
		return err
	}

	var flag model.ExecFlags
	flag.NoAutoRestart = false
	// remove -s
	flag.Args = strings.Join(os.Args[2:], " ")
	os.Setenv("PMON2_CONF", path+"/pmon2_self_conf.yml")
	ServiceRun([]string{cmdPath}, flag, path)
	go func() {
		// start monitor service
		service_monitor.NewMonitor()
	}()
	return nil
}
