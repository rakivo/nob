package nob

import (
	"fmt"
)

var (
	cs children
)

func (c *Child) Wait() error {
	if err := c.Cmd.Wait(); err != nil {
		fmt.Print(c.Stderr.String())
		return err
	}

	return nil
}

func (c *Child) MustWait() {
	if err := c.Wait(); err != nil {
		panic(err)
	}
}

func WaitAll() error {
	for len(cs.buf) > 0 {
		n := len(cs.buf)

		last := cs.buf[n-1]

		cs.mu.Lock()
		cs.buf = cs.buf[:n-1]
		cs.mu.Unlock()

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
