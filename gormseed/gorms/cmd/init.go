/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/raspiantoro/gormseeder/gormseed/internal/app"
	"github.com/spf13/cobra"
)

var (
	flagWithCli bool
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize gormseed",
	Long:  `Initialize gormseed will create new directory for your gormseed files`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := app.Config{
			GormseedDir: flagDir,
			WithCli:     flagWithCli,
		}
		app.Init(cfg)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.
	initCmd.Flags().StringVarP(&flagDir, "dir", "d", "", "gormseed files directory")
	initCmd.Flags().BoolVarP(&flagWithCli, "with-cli", "c", false, "create cli")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
