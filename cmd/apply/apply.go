package apply

import (
	"github.com/CausQ-Networks/praxis-cli/cmd/apply/blueprint"
	"github.com/spf13/cobra"
)

func ApplyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply commands",
		Long:  "Apply resources to the Praxis server",
	}

	cmd.AddCommand(blueprint.BlueprintCmd())

	return cmd
}
