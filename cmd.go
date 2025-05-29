package nob

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Cmd struct {
	buf []string
	env []string
}

func New(cmd ...string) *Cmd {
	return &Cmd{buf: cmd}
}

func (c *Cmd) Append(cmd ...string) {
	c.buf = append(c.buf, cmd...)
}

func (c *Cmd) Appendf(format string, a ...any) {
	s := fmt.Sprintf(format, a...)
	c.buf = append(c.buf, s)
}

func (c *Cmd) AppendEnv(v ...string) {
	c.env = append(c.env, v...)
}

func (c *Cmd) AppendfEnv(format string, a ...any) {
	s := fmt.Sprintf(format, a...)
	c.env = append(c.env, s)
}

func (c *Cmd) Run() (*Child, error) {
	rendered := strings.Join(c.buf, " ")
	fmt.Println(rendered)

	cmd := exec.Command(c.buf[0], c.buf[1:]...)
	cmd.Env = append(os.Environ(), c.env...)

	stdoutBuf := &bytes.Buffer{}
	stderrBuf := &bytes.Buffer{}

	cmd.Stdout = io.MultiWriter(stdoutBuf, os.Stdout) // tee stdout -> stdout
	cmd.Stderr = io.MultiWriter(stderrBuf, os.Stdout) // tee stderr -> stdout

	if err := cmd.Start(); err != nil {
		return &Child{Stderr: stderrBuf}, err
	}

	child := &Child{
		Cmd:    cmd,
		Stdout: stdoutBuf,
		Stderr: stderrBuf,
	}

	cs.append(child)

	return child, nil
}
