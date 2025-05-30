package nob

import (
	"context"
	"fmt"
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

func (s *Session) Pending() []*Process {
	return s.pending
}

// starts cmd under ctx, streams live output if configured,
// registers the resulting Process for later waiting.
func (s *Session) StartContext(ctx context.Context, c *Cmd) (*Process, error) {
	return s.startContext(ctx, c)
}

func (s *Session) Start(c *Cmd) (*Process, error) {
	return s.startContext(nil, c)
}

func (s *Session) MustStart(c *Cmd) *Process {
	p, err := s.Start(c)
	if err != nil {
		panic(err)
	}
	return p
}

// waits for every spawned process
func (s *Session) WaitAll() error {
	s.mu.Lock()
	procs := s.pending
	s.pending = nil
	s.mu.Unlock()

	var wg sync.WaitGroup
	errs := make(chan error, len(procs))

	for _, proc := range procs {
		wg.Add(1)
		go func(p *Process) {
			defer wg.Done()
			if err := p.cmd.Wait(); err != nil {
				errs <- err
			}
		}(proc)
	}

	wg.Wait()
	close(errs)

	// return the first error, if any
	for err := range errs {
		return err
	}

	return nil
}

// like WaitAll but panics on error.
func (s *Session) MustWaitAll() {
	if err := s.WaitAll(); err != nil {
		panic(err)
	}
}

func (s *Session) startContext(ctx context.Context, c *Cmd) (*Process, error) {
	var raw *exec.Cmd

	if ctx != nil {
		raw = c.RawContext(ctx)
	} else {
		raw = c.Raw()
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
