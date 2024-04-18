package bot

import (
	"air-quality-bot/internal/services/aqicn"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

const (
	commandStatus = "status"
)

type Bot struct {
	bot    *tgbotapi.BotAPI
	aqicn  *aqicn.Service
	logger *zap.Logger
}

func (b *Bot) Start() {
	b.logger.Info("starting bot")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			switch update.Message.Command() {
			case commandStatus:
				data := b.aqicn.Data()
				msgText := fmt.Sprintf(`
<b>Air Quality in Limassol [%d - %s]:</b>
- PM2.5: %0.2f
- PM10: %0.2f
- NO2: %0.2f
- CO: %0.2f
- SO2: %0.2f
- Ozone: %0.2f
- Primary pollutant: %s
- Humidity: %0.1f
- Pressure:  %0.1fmb
`,
					data.Data.Aqi,
					aqiValue(data.Data.Aqi),
					data.Data.Iaqi.Pm25.V,
					data.Data.Iaqi.Pm10.V,
					data.Data.Iaqi.No2.V,
					data.Data.Iaqi.Co.V,
					data.Data.Iaqi.So2.V,
					data.Data.Iaqi.O3.V,
					data.Data.Dominentpol,
					data.Data.Iaqi.H.V,
					data.Data.Iaqi.P.V,
				)

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
				msg.ParseMode = tgbotapi.ModeHTML

				_, err := b.bot.Send(msg)
				if err != nil {
					b.logger.Error("failed to send message", zap.Error(err))
				}
			}
		}
	}
}

func aqiValue(aqi int) string {
	switch {
	case aqi <= 50:
		return "Good"
	case aqi <= 100:
		return "Moderate"
	case aqi <= 150:
		return "Unhealthy for Sensitive Groups"
	case aqi <= 200:
		return "Unhealthy"
	case aqi <= 300:
		return "Very Unhealthy"
	}

	return "Hazardous"
}

func New(
	aqicn *aqicn.Service,
	logger *zap.Logger,
	config *Config,
) *Bot {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		logger.Fatal("Failed to connect to telegram bot", zap.Error(err))
	}

	bot.Debug = config.Debug

	return &Bot{
		bot:    bot,
		aqicn:  aqicn,
		logger: logger,
	}
}
