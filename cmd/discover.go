package cmd

import (
	"time"

	"github.com/pugkong/ylc/app"
	"github.com/pugkong/ylc/pokemon"
	"github.com/spf13/cobra"
)

var (
	discoverListen   *string
	discoverDuration *time.Duration
)

var discoverCmd = &cobra.Command{
	GroupID: manageGroup.ID,
	Use:     "discover",
	Aliases: []string{"d", "dis"},
	Short:   "Discover new or update knows bulbs",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		manager := app.NewManager(store, pokemon.NewNames(), cmd)

		return manager.Discover(*discoverListen, *discoverDuration)
	},
}

func init() {
	rootCmd.AddCommand(discoverCmd)

	discoverListen = discoverCmd.Flags().StringP("listen", "l", ":0", "address to listen")
	discoverDuration = discoverCmd.Flags().DurationP("duration", "d", time.Second, "time to listen")
}
