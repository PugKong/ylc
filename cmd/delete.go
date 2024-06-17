package cmd

import (
	"github.com/pugkong/ylc/app"
	"github.com/pugkong/ylc/pokemon"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	GroupID: manageGroup.ID,
	Use:     "delete [bulb name]",
	Aliases: []string{"del"},
	Short:   "Delete bulb",
	Args:    cobra.ExactArgs(1),
	ValidArgsFunction: func(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return store.AllNames(), cobra.ShellCompDirectiveDefault
		}

		return nil, cobra.ShellCompDirectiveDefault
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return app.NewManager(store, pokemon.NewNames(), cmd).Delete(args[0])
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
