package cmd

import (
	"MagicBox/utils"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use: "MagicBox",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				return
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(configCmd)
}

func Execute() error {
	err := rootCmd.Execute()
	if err != nil {
		utils.GLOBAL_LOGGER.Error("rootCmd error: " + err.Error())
	}
	return err
}
