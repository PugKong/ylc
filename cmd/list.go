package cmd

import (
	"github.com/pugkong/ylc/app"
	"github.com/pugkong/ylc/pokemon"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	GroupID: manageGroup.ID,
	Use:     "list",
	Aliases: []string{"l", "ls"},
	Short:   "List known bulbs",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		return app.NewManager(store, pokemon.NewNames(), cmd).List()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
