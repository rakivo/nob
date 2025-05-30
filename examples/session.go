package main

import (
	"fmt"

	"github.com/rakivo/nob"
)

func main() {
	session := nob.NewSession()

	// kick off a few long-running commands in parallel:
	cmd := nob.Command("sleep", "2")

	_ = session.MustStart(cmd)

	cmd = nob.Command("sleep", "1")

	_ = session.MustStart(cmd)

	fmt.Println("Both sleeps started, doing other work now…")

	// (do other work here)

	// finally, wait for **all** commands to finish:
	if err := session.WaitAll(); err != nil {
		fmt.Printf("one of the commands failed: %v\n", err)
	} else {
		fmt.Println("All done!")
	}
}
