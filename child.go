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
	cs.buf = append(cs.buf, c)
}
