package main

import (
	"github.com/rakivo/nob"
)

func main() {
	s := nob.NewSession()
	defer s.MustWaitAll()

	cmd := nob.Command("echo", "Hello, world!")
	_, err := s.Run(cmd)
	if err != nil {
		panic(err)
	}
}
