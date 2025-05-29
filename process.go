package nob

import "os/exec"

// wraps an exec.Cmd that’s already been started
type Process struct {
	cmd *exec.Cmd
}

// waits for the process -> returns any error
func (p *Process) Wait() error {
	return p.cmd.Wait()
}

// returns the underlying OS process ID
func (p *Process) PID() int {
	if p.cmd.Process != nil {
		return p.cmd.Process.Pid
	}
	return 0
}

