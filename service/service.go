package public_service

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	service "github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/app/model"
	"github.com/ntt360/pmon2/app/output"
)

func ServiceRun(args []string, flag model.ExecFlags, path string) {
	// get exec abs file path
	execPath, err := GetExecFile(args)
	if err != nil {
		service.Log.Error(err.Error())
		return
	}
	flag.Log = path + "/pmon2/log/" + filepath.Base(execPath) + ".log"
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
		log.Fatal(err)
	}

	var flag model.ExecFlags
	flag.NoAutoRestart = false
	// remove -s
	flag.Args = strings.Join(os.Args[2:], " ")
	os.Setenv("PMON2_CONF", path+"/pmon2_self_conf.yml")
	ServiceRun([]string{cmdPath}, flag, path)
	return nil
}
