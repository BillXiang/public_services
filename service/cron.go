package public_service

import (
	"os"
	"public_service/config"
	"public_service/daemon"
	"strconv"

	"github.com/robfig/cron/v3"
)

func Cron(jobCfg *config.Config) {
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
			pd.Wait()
		})
	}
	cron.Start()

	// fmt.Println("%v", cron.Entries())
	entries := cron.Entries()
	for _, entry := range entries {
		os.WriteFile("/var/log/public_service/cron.log", []byte(strconv.FormatInt(int64(entry.ID), 10)), 0644)
	}
}
