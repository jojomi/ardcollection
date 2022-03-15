package main

import (
	"fmt"
	"os"
)

func main() {
	rootCmd := getRootCmd()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
