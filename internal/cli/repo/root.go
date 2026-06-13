package repo

import (
	"github.com/slingshot-pipelines/cli/internal/cli/repo/component"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repo",
		Short: "Work with a repository",
	}

	cmd.AddCommand(component.NewCommand())

	return cmd
}
