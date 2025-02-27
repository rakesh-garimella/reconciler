package cmd

import (
	startCmd "github.com/kyma-incubator/reconciler/cmd/service/start"
	"github.com/kyma-incubator/reconciler/internal/cli"
	"github.com/spf13/cobra"
)

func NewCmd(o *cli.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service",
		Short: "Manage Kyma reconciler service",
		Long:  "Administrative CLI tool for the Kyma reconciler service",
	}

	cmd.AddCommand(startCmd.NewCmd(startCmd.NewOptions(o)))

	return cmd
}
