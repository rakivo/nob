package nob

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type command struct {
	name   string
	args   []string
	env    []string
	stdout io.Writer
	stderr io.Writer
}

func Command(name string, args ...string) *command {
	return &command{
		name:   name,
		args:   args,
		stdout: os.Stdout, // stdout -> stdout
		stderr: os.Stdout, // stderr -> stderr
	}
}

func (c *command) WithEnv(env ...string) *command {
	c.env = append(c.env, env...)
	return c
}

func (c *command) WithOutput(out, errOut io.Writer) *command {
	c.stdout = out
	c.stderr = errOut
	return c
}

func (c *command) Append(cmd ...string) {
	c.args = append(c.args, cmd...)
}

func (c *command) Appendf(format string, a ...any) {
	s := fmt.Sprintf(format, a...)
	c.args = append(c.args, s)
}

func (c *command) AppendEnv(v ...string) {
	c.env = append(c.env, v...)
}

func (c *command) AppendfEnv(format string, a ...any) {
	s := fmt.Sprintf(format, a...)
	c.env = append(c.env, s)
}

func (c *command) CombinedOutput() ([]byte, error) {
	cmd := exec.Command(c.name, c.args...)
	cmd.Env = append(os.Environ(), c.env...)
	return cmd.CombinedOutput()
}

func (c *command) Render() string {
	var n int

	n += len(c.name)

	for _, a := range c.args {
		n += 1 + len(a) // space + len
	}

	var b strings.Builder
	b.Grow(n)

	b.WriteString(c.name)

	for _, a := range c.args {
		b.WriteByte(' ')
		b.WriteString(a)
	}

	return b.String()
}
