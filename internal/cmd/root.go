package cmd

import (
	"os"

	"github.com/slingshot-pipelines/cli/internal/cli/component"
	"github.com/spf13/cobra"
)

func Execute() {
	root := &cobra.Command{
		Use:   "slingshot",
		Short: "CLI for working with slingshot pipelines",
	}

	root.AddCommand(component.NewCommand())

	err := root.Execute()
	if err != nil {
		os.Exit(1)
	}
}
