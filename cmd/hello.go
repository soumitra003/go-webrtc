package cmd

import (
	"github.com/soumitra003/goframework/config"
	"github.com/soumitra003/goframework/logging"
	"github.com/spf13/cobra"
)

func NewHelloCommand() *cobra.Command {

	cfg, err := config.Load("")
	if err != nil {
		panic(err)
	}
	logger := logging.GetLogger()

	helloCmd := &cobra.Command{
		Use:   "hello",
		Short: "Say hello",
		Long:  "This command is saying hello to you",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Sugar().Infof("Hellow from %s environment", cfg.Environment)
			return nil
		},
	}

	return helloCmd
}
