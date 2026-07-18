//go:build !cgo

package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "DevBox GUI requires CGO. Please build with CGO_ENABLED=1 and a working C compiler.")
	fmt.Fprintln(os.Stderr, "Use the CLI interface instead: devbox-cli")
	os.Exit(1)
}
