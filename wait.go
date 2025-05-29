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
	for c := cs.pop(); c != nil; c = cs.pop() {
		if err := c.Cmd.Wait(); err != nil {
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
