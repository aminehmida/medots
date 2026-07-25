/*
Copyright © 2022 Amine Hmida <aminehmida@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/aminehmida/medots/config"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a configuration for all or a specific application",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		configPath := cmd.Flag("config").Value.String()
		conf, err := config.ReadConfig(configPath)
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(args) == 0 {
			fmt.Println("No application name specified. Deploying all applications.")
			for app, entries := range *conf {
				deployApp(app, entries)
			}
		} else {
			for _, arg := range args {
				entries, ok := (*conf)[arg]
				if !ok {
					color.Red("No configuration found for: " + arg)
					continue
				}
				deployApp(arg, entries)
			}
		}
	},
}

// deployApp deploys all config entries for a single application.
func deployApp(app string, entries []config.AppConfig) {
	fmt.Println("Deploying config for: " + app)
	for _, entry := range entries {
		stdout, stderr, err := entry.Link()
		if err != nil {
			color.Red(err.Error())
		}
		if stdout != nil {
			color.Green("stdout: \n" + *stdout)
		}
		if stderr != nil {
			color.Red("stderr: \n" + *stderr)
		}
	}
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
