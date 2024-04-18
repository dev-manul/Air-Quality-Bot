package config

import (
	"air-quality-bot/internal/services/aqicn"
	"air-quality-bot/internal/services/bot"
	"go.uber.org/fx"
)

type Config struct {
	fx.Out `required:"false"`

	Bot   *bot.Config
	Aqicn *aqicn.Config
}
