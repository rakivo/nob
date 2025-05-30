package nob

import (
	"context"
	"fmt"
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

func (s *Session) Pending() []*Process {
	return s.pending
}

// starts cmd under ctx, streams live output if configured,
// registers the resulting Process for later waiting.
func (s *Session) StartContext(ctx context.Context, c *command) (*Process, error) {
	raw := c.Raw()

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

func (s *Session) Start(c *command) (*Process, error) {
	return s.StartContext(context.Background(), c)
}

func (s *Session) MustStart(c *command) *Process {
	p, err := s.Start(c)
	if err != nil {
		panic(err)
	}
	return p
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
