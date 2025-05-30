package nob

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Cmd struct {
	Name   string
	Args   []string
	Env    []string
	Stdout io.Writer
	Stderr io.Writer
}

func Command(name string, args ...string) *Cmd {
	return &Cmd{
		Name:   name,
		Args:   args,
		Env:    os.Environ(),
		Stdout: os.Stdout, // stdout -> stdout
		Stderr: os.Stdout, // stderr -> stderr
	}
}

func (c *Cmd) WithOutput(out, errOut io.Writer) *Cmd {
	c.Stdout = out
	c.Stderr = errOut
	return c
}

func (c *Cmd) WithArgs(cmd ...string) *Cmd {
	c.Args = append(c.Args, cmd...)
	return c
}

func (c *Cmd) WithArgsf(format string, a ...any) *Cmd {
	s := fmt.Sprintf(format, a...)
	c.Args = append(c.Args, s)
	return c
}

func (c *Cmd) WithEnv(v ...string) *Cmd {
	c.Env = append(c.Env, v...)
	return c
}

func (c *Cmd) WithEnvf(format string, a ...any) *Cmd {
	s := fmt.Sprintf(format, a...)
	c.Env = append(c.Env, s)
	return c
}

func (c *Cmd) Render() string {
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

func (c *Cmd) RawContext(ctx context.Context) *exec.Cmd {
	return c.raw(ctx)
}

func (c *Cmd) Raw() *exec.Cmd {
	return c.raw(nil)
}

func (c *Cmd) Run() error {
	return c.Raw().Run()
}

func (c *Cmd) MustRun() {
	if err := c.Raw().Run(); err != nil {
		panic(err)
	}
}

func (c *Cmd) Output() ([]byte, error) {
	cmd := c.Raw()
	return cmd.Output()
}

func (c *Cmd) CombinedOutput() ([]byte, error) {
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

func (c *Cmd) raw(ctx context.Context) *exec.Cmd {
	var raw *exec.Cmd

	if ctx != nil {
		raw = exec.CommandContext(ctx, c.Name, c.Args...)
	} else {
		raw = exec.Command(c.Name, c.Args...)
	}

	raw.Env = append(os.Environ(), c.Env...)

	if c.Stdout != nil {
		raw.Stdout = c.Stdout
	}

	if c.Stderr != nil {
		raw.Stderr = c.Stderr
	}

	return raw
}

