package bot

import (
	"fmt"
	"hellobox/database"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandlePartnerBot() {
	bot, err := tgbotapi.NewBotAPI("5009480809:AAHP6dwtinjlHTfKk-pHRlHWFd-_dbOPo3M")
	if err != nil {
		panic(err)
	}
	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)

	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.CallbackQuery != nil {
			s := strings.Split(update.CallbackQuery.Data, "#")
			switch s[0] {
			case "activate":
				id, _ := strconv.Atoi(s[1])
				product := database.GetCartProductsById(uint(id))
				product.Utilized = true
				database.EditCartProduct(*product)
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Success! Product has been utilized!")
				bot.Send(msg)
			}
			continue
		}
		switch update.Message.Text {
		case "/start":
			message := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome")
			bot.Send(message)
			continue
		}
		product := database.GetCartProductsByToken(update.Message.Text)
		if product != nil && !product.Utilized && product.Id != 0 {
			reply := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("✅ Aктивировать", fmt.Sprintf("activate#%d", product.ProductId))))
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Продукт: %s\nКоличество: %d", product.Product.Name, product.Count))
			msg.ReplyMarkup = reply
			bot.Send(msg)
		} else {
			message := tgbotapi.NewMessage(update.Message.Chat.ID, "The product already expired,utilized or does not exist")
			bot.Send(message)
		}
		// sendPartnerMenu()
	}
}
