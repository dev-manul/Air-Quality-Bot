package main

import (
	"air-quality-bot/internal/config"
	"air-quality-bot/internal/services/aqicn"
	"air-quality-bot/internal/services/bot"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"time"
)

func main() {
	fx.New(
		fx.Provide(zap.NewDevelopment),
		fx.Provide(config.Provide(config.Config{}, "config.dev.toml", "config.toml")),
		fx.Provide(
			aqicn.New,
			bot.New,
		),
		fx.Invoke(func(bot *bot.Bot) {
			go bot.Start()
		}),
		fx.Invoke(func(s *aqicn.Service) {
			go func() {
				for {
					s.Update()
					time.Sleep(30 * time.Minute)
				}
			}()
		}),
	).Run()
}
