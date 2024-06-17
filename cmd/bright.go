package cmd

import (
	"fmt"
	"strconv"

	"github.com/pugkong/ylc/app"
	"github.com/pugkong/ylc/yeelight"
	"github.com/spf13/cobra"
)

var (
	brightBackground *bool
	brightEffect     *yeelight.Effect = &yeelight.EffectSmooth
	brightDuration   *int
)

var brightCmd = &cobra.Command{
	GroupID: controlGroup.ID,
	Use:     "bright [bulb name] [bright]",
	Aliases: []string{"b"},
	Short:   "Set bright",
	Args:    cobra.ExactArgs(2),
	ValidArgsFunction: func(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return store.AllNames(), cobra.ShellCompDirectiveDefault
		}

		if len(args) == 1 {
			brights := make([]string, 0)
			for i := 10; i <= 100; i += 10 {
				brights = append(brights, strconv.Itoa(i))
			}

			return brights, cobra.ShellCompDirectiveDefault
		}

		return nil, cobra.ShellCompDirectiveDefault
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		value, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("parse bright: %w", err)
		}

		control := app.NewControl(store, cmd)

		if *brightBackground {
			return control.SetBackgroundBright(name, value, *brightEffect, *brightDuration)
		}

		return control.SetBright(name, value, *brightEffect, *brightDuration)
	},
}

func init() {
	rootCmd.AddCommand(brightCmd)

	brightBackground = brightCmd.Flags().Bool("bg", false, "set background bright")
	brightCmd.Flags().VarP(newEffectValue(brightEffect), "effect", "e", "smooth or sudden")
	brightDuration = brightCmd.Flags().IntP("duration", "d", 500, "effect duration (ms)")
}
