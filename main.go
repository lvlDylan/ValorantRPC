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
	_ "embed"
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

//go:embed ico.ico
var appIcon []byte

func onReady() {
	cfg, err := config.Load("config.json")
	if err != nil {
		panic(err)
	}

	logger, err := logging.New(cfg.Logging)
	if err != nil {
		panic(err)
	}

	translator, err := i18n.NewService(cfg.Language)
	if err != nil {
		logger.Error("Error initializing i18n service, fallback to English", "error", err)
	}

	logger.Info("Initializing Valorant RPC...")

	systray.SetTitle("ValorantRPC")
	systray.SetTooltip("Valorant Rich Presence")
	systray.SetIcon(appIcon)

	mStatus := systray.AddMenuItem(translator.T("app.menus.waiting"), "")
	mStatus.Disable()

	mScanNow := systray.AddMenuItem(translator.T("app.menus.scan_now"), translator.T("app.menus.scan_now.tooltip"))
	systray.AddSeparator()

	mLang := systray.AddMenuItem(translator.T("app.menus.language"), "")
	mEn := mLang.AddSubMenuItem("English", "")
	mFr := mLang.AddSubMenuItem("Français", "")
	mEn.Show()
	mFr.Show()
	systray.AddSeparator()

	mQuit := systray.AddMenuItem(translator.T("app.menus.quit"), "")

	riotClient := client.New()

	registry, err := assets.NewRegistry()
	if err != nil {
		logger.Error("Error initializing Valorant assets registry", "error", err)
		return
	}

	discordClient = discord.New(DiscordAppId, translator, registry, logger)
	appEngine = engine.New(riotClient, discordClient, logger, 15*time.Second)

	ctx, cancel := context.WithCancel(context.Background())
	triggerScan := make(chan struct{})
	statusChan := make(chan bool)

	go func() {
		for isRunning := range statusChan {
			if isRunning {
				mStatus.SetTitle(translator.T("app.menus.running"))
				mScanNow.Disable()
			} else {
				mStatus.SetTitle(translator.T("app.menus.waiting"))
				mScanNow.Enable()
			}
		}
	}()

	switchLanguage := func(newLang string) {
		logger.Info("Switching language", "language", newLang)

		cfg.Language = newLang
		if err := config.Save("config.json", cfg); err != nil {
			logger.Error("Error switching language", "error", err)
		}

		if err := translator.SetLang(newLang); err != nil {
			logger.Error("Error switching language", "error", err)
			return
		}

		mScanNow.SetTitle(translator.T("app.menus.scan_now"))
		mLang.SetTitle(translator.T("app.menus.language"))
		mQuit.SetTitle(translator.T("app.menus.quit"))

		select {
		case triggerScan <- struct{}{}:
		default:
		}
	}

	mEn.Click(func() { switchLanguage("en") })
	mFr.Click(func() { switchLanguage("fr") })

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
