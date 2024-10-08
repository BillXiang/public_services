package public_service

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/app/model"
	"github.com/ntt360/pmon2/client/proxy"
)

func LoadFirst(execPath string, flags string) ([]string, error) {
	data, err := proxy.RunProcess([]string{"start", execPath, flags})
	if err != nil {
		return nil, err
	}

	var tbData []string
	_ = json.Unmarshal(data, &tbData)

	return tbData, nil
}

// check the process already have
func ProcessExist(execPath string) (*model.Process, bool) {
	var process model.Process
	err := app.Db().First(&process, "process_file = ?", execPath).Error
	if err != nil {
		return nil, false
	}

	return &process, true
}

func GetExecFile(args []string) (string, error) {
	execFile := args[0]
	_, err := os.Stat(execFile)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("%s not exist", execFile)
	}

	if path.IsAbs(execFile) {
		return execFile, nil
	}

	absPath, err := filepath.Abs(execFile)
	if err != nil {
		return "", fmt.Errorf("get file path error: %v", err.Error())
	}

	return absPath, nil
}
