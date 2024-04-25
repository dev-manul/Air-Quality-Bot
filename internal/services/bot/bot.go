package bot

import (
	"air-quality-bot/internal/services/aqicn"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"strings"
	"time"
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
				forecast := []string{}
				for _, value := range data.Data.Forecast.Daily.Pm25 {
					day, err := time.Parse(value.Day, "2006-01-02")
					if err != nil {
						b.logger.Error("failed to parse date", zap.String("date", value.Day), zap.Error(err))
					}
					if day.After(time.Now()) {
						forecast = append(forecast, fmt.Sprintf("%s: %d", value.Day, value.Avg))
					}
				}
				msgText := fmt.Sprintf(`
<b>Air Quality in Limassol [%d - %s]:</b>
- PM2.5: %0.2f [Good < 50]
- PM10: %0.2f [Good < 50 ]
- NO2: %0.2f [Good < 50]
- CO: %0.2f [Good < 50]
- SO2: %0.2f [Good < 50]
- Ozone: %0.2f [Good < 50 ]
- Primary pollutant: %s
- Humidity: %0.1f
- Pressure:  %0.1fmb [Normal 1013.25mb]
- Forecast for PM2.5: %s
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
					strings.Join(forecast, ", "),
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
