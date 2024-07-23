package public_service

import (
	"public_service/config"
	"public_service/daemon"
)

func Exec(jobCfg *config.Config) {
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
		pd.Wait()
	}
}
