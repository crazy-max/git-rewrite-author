package main

import (
	"fmt"
	"os"

	"github.com/crazy-max/git-rewrite-author/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
