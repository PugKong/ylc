package cmd

import (
	"fmt"
	"strconv"

	"github.com/pugkong/ylc/app"
	"github.com/pugkong/ylc/yeelight"
	"github.com/spf13/cobra"
)

var (
	temperatureBackground *bool
	temperatureEffect     = &yeelight.EffectSmooth
	temperatureDuration   *int
)

var temperatureCmd = &cobra.Command{
	GroupID: controlGroup.ID,
	Use:     "temperature [bulb name] [temperature]",
	Aliases: []string{"t", "temp"},
	Short:   "Set color temperature",
	Args:    cobra.ExactArgs(2),
	ValidArgsFunction: func(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return store.AllNames(), cobra.ShellCompDirectiveDefault
		}

		if len(args) == 1 {
			temperatures := make([]string, 0)
			for i := 1700; i <= 6500; i += 100 {
				temperatures = append(temperatures, strconv.Itoa(i))
			}

			return temperatures, cobra.ShellCompDirectiveDefault
		}

		return nil, cobra.ShellCompDirectiveDefault
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		value, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("parse temperature: %w", err)
		}

		control := app.NewControl(store, cmd)

		if *temperatureBackground {
			return control.SetBackgroundTemperature(name, value, *temperatureEffect, *temperatureDuration)
		}

		return control.SetTemperature(name, value, *temperatureEffect, *temperatureDuration)
	},
}

func init() {
	rootCmd.AddCommand(temperatureCmd)

	temperatureBackground = temperatureCmd.Flags().Bool("bg", false, "set background temperature")
	temperatureCmd.Flags().VarP(newEffectValue(temperatureEffect), "effect", "e", "smooth or sudden")
	temperatureDuration = temperatureCmd.Flags().IntP("duration", "d", 500, "effect duration (ms)")
}
