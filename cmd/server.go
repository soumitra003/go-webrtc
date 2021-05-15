package cmd

import (
	"context"
	"fmt"

	"github.com/soumitra003/go-webrtc/auth"
	"github.com/soumitra003/go-webrtc/sse"
	"github.com/soumitra003/goframework/config"
	"github.com/soumitra003/goframework/example/hello"
	"github.com/soumitra003/goframework/logging"
	"github.com/soumitra003/goframework/server"
	"github.com/spf13/cobra"
)

func NewServerCommand(ctx context.Context) *cobra.Command {

	serverCommand := &cobra.Command{
		Use:   "server",
		Short: "Start server",
		RunE: func(cmd *cobra.Command, args []string) error {

			cfg, err := config.Load("")
			if err != nil {
				panic(err)
			}
			logger := logging.GetLogger()

			fmt.Print(ctx, "\n")
			logger.Sugar().Debugf("Environment: %s", cfg.Environment)
			_ = RunServerStart(ctx, cfg)
			return nil
		},
	}

	return serverCommand
}

func RunServerStart(ctx context.Context, cfg *config.Config) error {

	svr := server.New(
		server.WithGlobalConfig(cfg),
	)

	svr.AddModule("hello", hello.New(*cfg))
	svr.AddModule("auth", auth.New(*cfg))
	svr.AddModule("sse", sse.New(*cfg))

	_ = svr.Start(ctx)
	logging.GetLogger().Info("Shutting down")

	return nil
}
