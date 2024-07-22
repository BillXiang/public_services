package public_service

import (
	"encoding/json"

	"github.com/ntt360/pmon2/app/model"
	"github.com/ntt360/pmon2/client/proxy"
)

func Restart(m *model.Process, flags string) ([]string, error) {
	newData, err := ReloadProcess(m, flags)
	if err != nil {
		return nil, err
	}

	return newData, nil
}

func ReloadProcess(m *model.Process, flags string) ([]string, error) {
	data, err := proxy.RunProcess([]string{"restart", m.ProcessFile, flags})

	if err != nil {
		return nil, err
	}

	var rel []string
	err = json.Unmarshal(data, &rel)
	if err != nil {
		return nil, err
	}

	return rel, nil
}
