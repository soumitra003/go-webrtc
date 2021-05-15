package cmd

import (
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var (
	rootCmd *cobra.Command
)

func InitRootCommand() error {
	ctx := context.Background()
	rootCmd = &cobra.Command{
		Use:   "cobra",
		Short: "A generator for Cobra based Applications",
		Long:  "Cobra is a CLI library for Go that empowers applications.\nThis application is a tool to generate the needed files\nto quickly create a Cobra application.",
	}

	rootCmd.AddCommand(NewHelloCommand())
	rootCmd.AddCommand(NewServerCommand(ctx))
	return nil
}

func Execute() error {
	return rootCmd.Execute()
}
