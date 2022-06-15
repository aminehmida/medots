/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

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
		fmt.Println("deploy called")
		configPath := cmd.Flag("config").Value.String()
		println(configPath)
		conf, err := config.ReadConfig(configPath)
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(args) == 0 {
			fmt.Println("No application name specified. Deploying all applications.")
			for k, e := range *conf {
				fmt.Println("Deploying config for: " + k)
				for _, v := range e {
					stdout, stderr, err := v.Link()
					if err != nil {
						fmt.Println(err)
					}
					if stdout != nil {
						color.Green("stdout: \n" + *stdout)
					}
					if stderr != nil {
						color.Red("stderr: \n" + *stderr)
					}
				}
			}
		} else {
			for _, arg := range args {
				//TODO: check if arg is in config
				fmt.Println("Deploying config for: " + arg)
				for _, v := range (*conf)[arg] {
					stdout, stderr, err := v.Link()
					if err != nil {
						fmt.Println(err)
					}
					if stdout != nil {
						color.Green("stdout: \n" + *stdout)
					}
					if stderr != nil {
						color.Red("stderr: \n" + *stderr)
					}
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
