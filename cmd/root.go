package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "bakfietscalculator",
	Short: "Calculcate prices",
	Run: func(cmd *cobra.Command, _ []string) {
		cmd.Help()
	},
}

func Execute() {
	rootCmd.Execute()
}