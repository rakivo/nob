# Go Command Runner

> A lightweight Go package for managing and orchestrating external commands, inspired by [nob.h](https://github.com/tsoding/nob.h), but reimplemented in Go.

# Quick start

```console
$ go get github.com/rakivo/nob
```

# Usage
```go
package main

import (
	"github.com/rakivo/nob"
)

func main() {
	// Wait for all forked children
	defer nob.WaitAll()

	cmd := nob.New("echo", "Hello, world!")
	_, err := cmd.Run()
	if err != nil {
			panic(err)
	}
}
```
