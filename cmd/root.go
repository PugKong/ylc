package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/pugkong/ylc/app"
	"github.com/spf13/cobra"
)

var (
	manageGroup  = cobra.Group{ID: "manage", Title: "Bulb Management"}
	controlGroup = cobra.Group{ID: "control", Title: "Bulb Control"}
)

var store *app.BulbFileStore

var rootCmd = &cobra.Command{
	Use:   "ylc",
	Short: "A CLI tool to control your Yeelight bulbs",
	PersistentPreRunE: func(*cobra.Command, []string) error {
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			return fmt.Errorf("get user cache dir: %w", err)
		}

		appCacheDir := path.Join(cacheDir, "ylc")
		if err := os.MkdirAll(appCacheDir, 0o700); err != nil {
			return fmt.Errorf("make cache dir: %w", err)
		}

		store = app.NewBulbFileStore(appCacheDir)
		if err := store.Init(); err != nil {
			return fmt.Errorf("init bulb store: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddGroup(&manageGroup, &controlGroup)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
