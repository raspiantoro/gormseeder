/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	flagPathAddCmd string
)

var longDesc = `Create and add a new migration file to your Golang project. 
The seed name will be prefixed with the current datetime in the format 'YYYYMMDDhhmm' as a seed version/identifier in the generated file name.`

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [seed name]",
	Short: "Add new migration file.",
	Long:  longDesc,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Print("Error: add requires a name for the seed.\n\n")
			cmd.Help()
			return
		} else if len(args) > 1 {
			fmt.Print("Error: add only accept one arguments [seed name]\n\n")
			cmd.Help()
			return
		}

		fmt.Println("run")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&flagPathAddCmd, "dir", "d", "", "seed files directory")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
