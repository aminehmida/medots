/*
Copyright © 2022 Amine Hmida <aminehmida@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "medots",
	Short: "Cross-platform dot config files manager using a yaml file",
	Long: `medots is a cross-platform config files manager.

It uses a dots.yaml file to symlink your dot files (including OS-specific
ones) and run commands before and/or after symlinking, so you can deploy
your configuration quickly and consistently across Linux, macOS and WSL.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.medots.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().StringP("config", "c", "./dots.yaml", "Path to your config file")
}
