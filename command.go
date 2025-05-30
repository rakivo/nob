package nob

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type command struct {
	Name   string
	Args   []string
	Env    []string
	Stdout io.Writer
	Stderr io.Writer
}

func Command(name string, args ...string) *command {
	return &command{
		Name:   name,
		Args:   args,
		Env:    os.Environ(),
		Stdout: os.Stdout, // stdout -> stdout
		Stderr: os.Stdout, // stderr -> stderr
	}
}

func (c *command) WithOutput(out, errOut io.Writer) *command {
	c.Stdout = out
	c.Stderr = errOut
	return c
}

func (c *command) WithArgs(cmd ...string) *command {
	c.Args = append(c.Args, cmd...)
	return c
}

func (c *command) WithArgsf(format string, a ...any) *command {
	s := fmt.Sprintf(format, a...)
	c.Args = append(c.Args, s)
	return c
}

func (c *command) WithEnv(v ...string) *command {
	c.Env = append(c.Env, v...)
	return c
}

func (c *command) WithEnvf(format string, a ...any) *command {
	s := fmt.Sprintf(format, a...)
	c.Env = append(c.Env, s)
	return c
}

func (c *command) Render() string {
	var n int

	n += len(c.Name)

	for _, a := range c.Args {
		n += 1 + len(a) // space + len
	}

	var b strings.Builder
	b.Grow(n)

	b.WriteString(c.Name)

	for _, a := range c.Args {
		b.WriteByte(' ')
		b.WriteString(a)
	}

	return b.String()
}

func (c *command) Raw() *exec.Cmd {
	raw := exec.Command(c.Name, c.Args...)
	raw.Env = append(os.Environ(), c.Env...)

	if c.Stdout != nil {
		raw.Stdout = c.Stdout
	}

	if c.Stderr != nil {
		raw.Stderr = c.Stderr
	}

	return raw
}

func (c *command) Run() error {
	return c.Raw().Run()
}

func (c *command) MustRun() {
	if err := c.Raw().Run(); err != nil {
		panic(err)
	}
}

func (c *command) Output() ([]byte, error) {
	cmd := c.Raw()
	return cmd.Output()
}

func (c *command) CombinedOutput() ([]byte, error) {
	cmd := c.Raw()
	cmd.Stdout = nil; cmd.Stderr = nil;
	return cmd.CombinedOutput()
}

func Run(name string, args ...string) error {
	return Command(name, args...).Raw().Run()
}

func MustRun(name string, args ...string) {
	if err := Run(name, args...); err != nil {
		panic(err)
	}
}

func Output(name string, args ...string) ([]byte, error) {
	return Command(name, args...).Raw().Output()
}

func CombinedOutput(name string, args ...string) ([]byte, error) {
	return Command(name, args...).Raw().CombinedOutput()
}

