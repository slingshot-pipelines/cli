package component

import "github.com/spf13/cobra"

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "component",
		Short: "Work with components",
	}

	cmd.AddCommand(newListCommand())

	return cmd
}
