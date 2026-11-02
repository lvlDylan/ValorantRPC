package main

import (
	"ValorantRPC/pkg/assets"
	"ValorantRPC/pkg/client"
	"ValorantRPC/pkg/discord"
	"ValorantRPC/pkg/engine"
	"ValorantRPC/pkg/i18n"
	"ValorantRPC/pkg/logging"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const DiscordAppId = "1564597359281832076"

func main() {
	logger, err := logging.New(logging.Config{
		Level:      "DEBUG",
		JSONFormat: false,
		FilePath:   "",
	})
	if err != nil {
		panic(err)
	}

	logger.Info("Initializing Valorant RPC...")

	translator, err := i18n.NewService("fr")
	if err != nil {
		logger.Error("Error initializing i18n service, fallback to English", "error", err)
	}

	logger.Info("Initializing Valorant assets registry..")
	registry, err := assets.NewRegistry()
	if err != nil {
		logger.Error("Error initializing Valorant assets registry", "error", err)
		return
	}

	riotClient, err := client.New()
	if err != nil {
		panic(err)
	}

	discordClient := discord.New(DiscordAppId, translator, registry, logger)
	appEngine := engine.New(riotClient, discordClient, logger, 15*time.Second)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		appEngine.Start()
	}()

	<-stop
	logger.Info("Shutting down...")

	discordClient.Disconnect()
	logger.Info("ValorantRPC stopped.")
}
