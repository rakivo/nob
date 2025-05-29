package nob

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

var (
	children []*Child
)

type Child struct {
	Cmd    *exec.Cmd
	Stdout *bytes.Buffer
	Stderr *bytes.Buffer
}

type Command struct {
	buf []string
	env []string
}

func New(cmd ...string) *Command {
	return &Command{buf: cmd}
}

func (c *Command) Append(cmd ...string) {
	c.buf = append(c.buf, cmd...)
}

func (c *Command) Appendf(format string, a ...any) {
	s := fmt.Sprintf(format, a...)
	c.buf = append(c.buf, s)
}

func (c *Command) AppendEnv(v ...string) {
	c.env = append(c.env, v...)
}

func (c *Command) AppendfEnv(format string, a ...any) {
	s := fmt.Sprintf(format, a...)
	c.env = append(c.env, s)
}

func (c *Command) Run() (*Child, error) {
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

	child := &Child{Cmd: cmd, Stdout: stdoutBuf, Stderr: stderrBuf}
	children = append(children, child)

	return child, nil
}

func Wait(child *Child) error {
	if err := child.Cmd.Wait(); err != nil {
		fmt.Print(child.Stderr.String())
		return err
	}

	return nil
}

func MustWait(child *Child) {
	if err := Wait(child); err != nil {
		panic(err)
	}
}

func WaitAll() error {
	for len(children) > 0 {
		last := children[len(children)-1]
		children = children[:len(children)-1]

		if err := last.Cmd.Wait(); err != nil {
			return err
		}
	}

	return nil
}

func MustWaitAll() {
	if err := WaitAll(); err != nil {
		panic(err)
	}
}
