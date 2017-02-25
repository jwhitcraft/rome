package main

import (
	"github.com/jwhitcraft/rome/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	romeCmd := cmd.RootCmd
	doc.GenMarkdownTree(romeCmd, "./docs")
}
