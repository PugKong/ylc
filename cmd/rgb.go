package cmd

import (
	"fmt"
	"strconv"

	"github.com/pugkong/ylc/app"
	"github.com/pugkong/ylc/yeelight"
	"github.com/spf13/cobra"
)

var (
	rgbBackground *bool
	rgbEffect     *yeelight.Effect = &yeelight.EffectSmooth
	rgbDuration   *int
)

var rgbCmd = &cobra.Command{
	GroupID: controlGroup.ID,
	Use:     "rgb [bulb name] [color]",
	Short:   "Set RGB color",
	Args:    cobra.ExactArgs(2),
	ValidArgsFunction: func(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return store.AllNames(), cobra.ShellCompDirectiveDefault
		}

		return nil, cobra.ShellCompDirectiveDefault
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		value, err := strconv.ParseInt(args[1], 16, 32)
		if err != nil {
			return fmt.Errorf("parse color: %w", err)
		}

		control := app.NewControl(store, cmd)

		if *rgbBackground {
			return control.SetBackgroundRGB(name, int(value), *rgbEffect, *rgbDuration)
		}

		return control.SetRGB(name, int(value), *rgbEffect, *rgbDuration)
	},
}

func init() {
	rootCmd.AddCommand(rgbCmd)

	rgbBackground = rgbCmd.Flags().Bool("bg", false, "set background color")
	rgbCmd.Flags().VarP(newEffectValue(rgbEffect), "effect", "e", "smooth or sudden")
	rgbDuration = rgbCmd.Flags().IntP("duration", "d", 500, "effect duration (ms)")
}
