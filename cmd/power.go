package cmd

import (
	"github.com/pugkong/ylc/app"
	"github.com/spf13/cobra"
)

var powerBackground *bool

var powerCmd = &cobra.Command{
	GroupID: controlGroup.ID,
	Use:     "power [bulb name]",
	Aliases: []string{"p"},
	Short:   "Toggle bulb power",
	Args:    cobra.ExactArgs(1),
	ValidArgsFunction: func(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return store.AllNames(), cobra.ShellCompDirectiveDefault
		}

		return nil, cobra.ShellCompDirectiveDefault
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		control := app.NewControl(store, cmd)

		if *powerBackground {
			return control.BackgroundToggle(args[0])
		}

		return control.PowerToggle(args[0])
	},
}

func init() {
	rootCmd.AddCommand(powerCmd)

	powerBackground = powerCmd.Flags().Bool("bg", false, "toggle only background power")
}
