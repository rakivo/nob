package main

import (
	"fmt"
	"github.com/rakivo/nob"
)

func main() {
	// this will run “ls -l”, wait for it to finish,
	// and print any error returned

	// by default,
	// stdout and stderr of the command are piped into os.Stdout
	if err := nob.Run("ls", "-l"); err != nil {
		fmt.Printf("ls failed: %v\n", err)
	}

	// same thing, but capturing output manually:
	out, err := nob.Command("echo", "hello, world").CombinedOutput()
	if err != nil {
		fmt.Printf("echo failed: %v\n", err)
	} else {
		fmt.Printf("echo said: %s", out)
	}
}
