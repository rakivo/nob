package nob

import (
	"bytes"
	"os/exec"
	"sync"
)

type Child struct {
	Cmd    *exec.Cmd
	Stdout *bytes.Buffer
	Stderr *bytes.Buffer
}

type children struct {
	buf []*Child
	mu  sync.Mutex
}

func (cs *children) append(c *Child) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	cs.buf = append(cs.buf, c)
}

func (cs *children) pop() *Child {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	n := len(cs.buf)

	if n == 0 {
		return nil
	}

	last := cs.buf[n-1]

	cs.buf = cs.buf[:n-1]

	return last
}
