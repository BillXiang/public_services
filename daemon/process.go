package daemon

import (
	"os"
	"os/exec"
	"time"

	psutilProcess "github.com/shirou/gopsutil/process"
)

type ProcessDaemon struct {
	CmdPath string
	Args    []string
	cmd     *exec.Cmd
	proc    *psutilProcess.Process
}

func (pd *ProcessDaemon) StartSubprocess() error {
	cmd := exec.Command(pd.CmdPath, pd.Args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()

	if err == nil {
		pd.cmd = cmd
		pd.proc, _ = psutilProcess.NewProcess(int32(cmd.Process.Pid))
	}
	return err
}

// kill -9, exit with "signal: killed"
func (pd *ProcessDaemon) StopSubprocess() error {
	return pd.cmd.Process.Kill()
}

func (pd *ProcessDaemon) Wait() error {
	return pd.cmd.Wait()
}

func (pd *ProcessDaemon) CPUPercent() (float64, error) {
	return pd.proc.Percent(time.Second)
}

func (pd *ProcessDaemon) Pid() int {
	return pd.cmd.Process.Pid
}
