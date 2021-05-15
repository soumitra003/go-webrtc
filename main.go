package main

import (
	"github.com/soumitra003/go-webrtc/cmd"
	"github.com/soumitra003/goframework/logging"
)

func main() {

	logger := logging.GetLogger()
	initRootCommand()

	if err := cmd.Execute(); err != nil {
		logger.Error("Error occurred while executing command")
	}
}

func initRootCommand() {
	if err := cmd.InitRootCommand(); err != nil {
		panic(err)
	}
}
