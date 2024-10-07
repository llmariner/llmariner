package main

import (
	"github.com/llmariner/llmariner/cli/internal/root"
	"github.com/spf13/cobra/doc"
)

func main() {
	cmd := root.Cmd()
	cmd.DisableAutoGenTag = true
	if err := doc.GenReSTTree(cmd, "output"); err != nil {
		panic(err)
	}
}
