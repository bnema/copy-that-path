package main

import (
	"fmt"
	"os"

	"github.com/bnema/yank-that/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "yt: %v\n", err)
		os.Exit(1)
	}
}
