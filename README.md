# Go Command Runner

> A lightweight Go package for managing and orchestrating external commands, inspired by [nob.h](https://github.com/tsoding/nob.h), but reimplemented in Go.

# Quick start

```console
$ go get github.com/rakivo/nob
```

# Usage
> One of the examples from `examples` directory
```go
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
	out, err := nob.CombinedOutput("echo", "hello, world")
	if err != nil {
		fmt.Printf("echo failed: %v\n", err)
	} else {
		fmt.Printf("echo said: %s", out)
	}
}
```
