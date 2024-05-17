package main

import (
	"github.com/llm-operator/cli/internal/root"
	"github.com/spf13/cobra/doc"
)

func main() {
	cmd := root.Cmd()
	cmd.DisableAutoGenTag = true
	if err := doc.GenReSTTree(cmd, "output"); err != nil {
		panic(err)
	}
}
