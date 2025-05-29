package main

import (
	"github.com/rakivo/nob"
)

func main() {
	// wait for all forked children
	defer nob.MustWaitAll()

	cmd := nob.New("echo", "Hello, world!")
	_, err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
