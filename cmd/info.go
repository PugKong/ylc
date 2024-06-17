package cmd

import (
	"github.com/pugkong/ylc/app"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	GroupID: controlGroup.ID,
	Use:     "info [BULB NAME]",
	Aliases: []string{"i"},
	Short:   "Show bulb info",
	Args:    cobra.ExactArgs(1),
	ValidArgsFunction: func(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return store.AllNames(), cobra.ShellCompDirectiveDefault
		}

		return nil, cobra.ShellCompDirectiveDefault
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return app.NewControl(store, cmd).Info(args[0])
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
