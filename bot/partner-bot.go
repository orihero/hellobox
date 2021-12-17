package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func HandlePartnerBot() {
	bot, err := tgbotapi.NewBotAPI("2105220707:AAEdCoI3Zbd0MY0UHYnjB5qeesHBVOSN1UU")
	if err != nil {
		panic(err)
	}
	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)

	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)
	for update := range updates {

	}
}
