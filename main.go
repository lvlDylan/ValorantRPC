package main

import (
	"ValorantRPC/pkg/assets"
	"ValorantRPC/pkg/client"
	"ValorantRPC/pkg/config"
	"ValorantRPC/pkg/discord"
	"ValorantRPC/pkg/engine"
	"ValorantRPC/pkg/i18n"
	"ValorantRPC/pkg/logging"
	"context"
	"time"

	"github.com/energye/systray"
)

const DiscordAppId = "1564597359281832076"

var (
	discordClient *discord.Client
	appEngine     *engine.Engine
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	cfg, err := config.Load("config.json")
	if err != nil {
		panic(err)
	}

	logger, err := logging.New(cfg.Logging)
	if err != nil {
		panic(err)
	}

	logger.Info("Initializing Valorant RPC...")

	systray.SetTitle("ValorantRPC")
	systray.SetTooltip("Valorant Rich Presence")

	mStatus := systray.AddMenuItem("Status: Waiting for valorant", "")
	mStatus.Disable()

	mScanNow := systray.AddMenuItem("Scan Now", "Scan Now to find Valorant instance")
	systray.AddSeparator()

	mQuit := systray.AddMenuItem("Quit", "Quit ValorantRPC")

	translator, err := i18n.NewService("fr")
	if err != nil {
		logger.Error("Error initializing i18n service, fallback to English", "error", err)
	}

	registry, err := assets.NewRegistry()
	if err != nil {
		logger.Error("Error initializing Valorant assets registry", "error", err)
		return
	}

	riotClient, err := client.New()
	if err != nil {
		panic(err)
	}

	discordClient = discord.New(DiscordAppId, translator, registry, logger)
	appEngine = engine.New(riotClient, discordClient, logger, 15*time.Second)

	ctx, cancel := context.WithCancel(context.Background())
	triggerScan := make(chan struct{})
	statusChan := make(chan bool)

	go func() {
		for isRunning := range statusChan {
			if isRunning {
				mStatus.SetTitle("Status: Valorant Active")
			} else {
				mStatus.SetTitle("Status: Valorant not found")
			}
		}
	}()

	mScanNow.Click(func() {
		logger.Info("Starting manual valorant scan...")
		select {
		case triggerScan <- struct{}{}:
		default:

		}
	})

	mQuit.Click(func() {
		logger.Info("Shutting down application")
		cancel()
		systray.Quit()
	})

	go appEngine.Start(ctx, triggerScan, statusChan)
}

func onExit() {
	if discordClient != nil {
		discordClient.Disconnect()
	}
}
