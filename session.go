package nob

import (
	"fmt"
	"context"
	"os"
	"os/exec"
	"sync"
)

// tracks spawned processes and waits on them.
type Session struct {
	mu      sync.Mutex
	pending []*Process
}

// NewSession creates an empty Session.
func NewSession() *Session {
	return &Session{}
}

// starts cmd under ctx, streams live output if configured,
// registers the resulting Process for later waiting.
func (s *Session) RunContext(ctx context.Context, c *command) (*Process, error) {
	raw := exec.CommandContext(ctx, c.name, c.args...)

	raw.Env = append(os.Environ(), c.env...)

	if c.stdout != nil {
		raw.Stdout = c.stdout
	}

	if c.stderr != nil {
		raw.Stderr = c.stderr
	}

	proc := &Process{cmd: raw}

	s.mu.Lock()
	s.pending = append(s.pending, proc)
	s.mu.Unlock()

	render := c.Render()
	fmt.Println(render)

	if err := raw.Start(); err != nil {
		return nil, err
	}

	return proc, nil
}

func (s *Session) Run(c *command) (*Process, error) {
	return s.RunContext(context.Background(), c)
}

// waits for every spawned process in LIFO order.
func (s *Session) WaitAll() error {
	for {
		s.mu.Lock()
		n := len(s.pending)
		if n == 0 {
			s.mu.Unlock()
			break
		}

		proc := s.pending[n-1]
		s.pending = s.pending[:n-1]
		s.mu.Unlock()

		if err := proc.cmd.Wait(); err != nil {
			return err
		}
	}
	return nil
}

// WaitAll but panics on error.
func (s *Session) MustWaitAll() {
	if err := s.WaitAll(); err != nil {
		panic(err)
	}
}

