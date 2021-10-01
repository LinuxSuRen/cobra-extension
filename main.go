package main

import (
	"github.com/linuxsuren/cobra-extension/pkg"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "cobra-extension",
		Short: "This is a demo command for cobra-extension",
	}

	root.AddCommand(pkg.NewCompletionCmd(root))

	if err := root.Execute(); err != nil {
		panic(err)
	}
}
