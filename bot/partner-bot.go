package bot

import (
	"fmt"
	"hellobox/database"
	"hellobox/env"
	"hellobox/models"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type PartnerCounts struct {
	Count     uint
	ProductId uint
}

var counts map[int64]PartnerCounts = make(map[int64]PartnerCounts)

func incrementActiveCount(userId int64, max int) bool {
	if val, ok := counts[userId]; ok && val.Count+1 <= uint(max) {
		counts[userId] = PartnerCounts{
			ProductId: val.ProductId,
			Count:     val.Count + 1,
		}
		return true
	} else {
		return false
	}
}

func decrementActiveCount(userId int64, max int) bool {
	if val, ok := counts[userId]; ok && val.Count-1 > 0 {
		counts[userId] = PartnerCounts{
			ProductId: val.ProductId,
			Count:     val.Count - 1,
		}
		return true
	} else {
		return false
	}
}

func checkProduct(chatId int64, userId int64, text string) {
	product := database.GetCartProductsByToken(text)
	if product != nil && product.Utilized == 0 && product.Id != 0 {
		user := database.FilterUser(models.User{ChatId: userId})
		if product.Product.PartnerId != user.PartnerId {
			message := tgbotapi.NewMessage(userId, "The product does not belong to you")
			env.PartnerBot.Send(message)
			return
		}
		count := 0
		if val, ok := counts[userId]; ok {
			if product.ProductId == val.ProductId {
				count = int(val.Count)
			} else {
				// When user picks new product
				counts[userId] = PartnerCounts{Count: uint(count), ProductId: product.ProductId}
			}
		} else {
			// When user types new product
			counts[userId] = PartnerCounts{Count: 0, ProductId: product.ProductId}
		}
		reply := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("➖", fmt.Sprintf("minus#%d", product.Id)),
				tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(count), "none"),
				tgbotapi.NewInlineKeyboardButtonData("➕", fmt.Sprintf("plus#%d", product.Id)),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("✅ Aктивировать", fmt.Sprintf("activate#%d", product.Id))))
		msg := tgbotapi.NewMessage(chatId, fmt.Sprintf("Продукт: %s\nКоличество: %d", product.Product.Name, product.Count))
		msg.ReplyMarkup = reply
		env.PartnerBot.Send(msg)
	} else {
		message := tgbotapi.NewMessage(chatId, "The product already expired,utilized or does not exist")
		env.PartnerBot.Send(message)
	}
}

func HandlePartnerBot() {
	bot, err := tgbotapi.NewBotAPI("5009480809:AAHP6dwtinjlHTfKk-pHRlHWFd-_dbOPo3M")
	if err != nil {
		panic(err)
	}
	bot.Debug = true
	env.PartnerBot = bot
	updateConfig := tgbotapi.NewUpdate(0)

	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)
	for update := range updates {
		if hasContact(update, true) {
			continue
		}
		if update.CallbackQuery != nil {
			user := database.FilterUser(models.User{ChatId: update.CallbackQuery.From.ID})
			if user.Id == 0 || user.PartnerId == 0 {
				//TODO check if the product on this partner
				env.Bot.Send(tgbotapi.NewCallbackWithAlert(update.CallbackQuery.ID, "You are not authorized to use this bot"))
				continue
			}
			s := strings.Split(update.CallbackQuery.Data, "#")
			switch s[0] {
			case "activate":
				id, _ := strconv.Atoi(s[1])
				product := database.GetCartProductsById(uint(id))
				if product.Product.PartnerId != user.PartnerId {
					env.Bot.Send(tgbotapi.NewCallbackWithAlert(update.CallbackQuery.ID, "You are not authorized to utilize this product"))
					continue
				}
				userId := update.CallbackQuery.From.ID
				if val, ok := counts[userId]; ok {
					if val.Count == product.Count {
						product.Utilized = 1
						product.Count = 0
					} else {
						product.Count = product.Count - val.Count
					}
				}
				database.EditCartProduct(*product)
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Success! Product has been utilized!")
				bot.Send(msg)
			case "plus":
				id, _ := strconv.Atoi(s[1])
				product := database.GetCartProductsById(uint(id))
				if product.Product.PartnerId != user.PartnerId {
					env.Bot.Send(tgbotapi.NewCallbackWithAlert(update.CallbackQuery.ID, "You are not authorized to utilize this product"))
					continue
				}
				ok := incrementActiveCount(update.CallbackQuery.From.ID, int(product.Count))
				if ok {
					checkProduct(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.From.ID, product.Token)
					bot.Send(tgbotapi.NewDeleteMessage(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID))
				}
			case "minus":
				id, _ := strconv.Atoi(s[1])
				product := database.GetCartProductsById(uint(id))
				if product.Product.PartnerId != user.PartnerId {
					env.Bot.Send(tgbotapi.NewCallbackWithAlert(update.CallbackQuery.ID, "You are not authorized to utilize this product"))
					continue
				}
				ok := decrementActiveCount(update.CallbackQuery.From.ID, int(product.Count))
				if ok {
					checkProduct(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.From.ID, product.Token)
					bot.Send(tgbotapi.NewDeleteMessage(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID))
				}
			}
			continue
		}
		if update.Message == nil || update.Message.Text == "" {
			continue
		}
		switch update.Message.Text {

		case "/start":
			reply := tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButtonContact("Отправить ваши контакы"),
				),
			)
			message := tgbotapi.NewMessage(update.Message.Chat.ID, "Чтобы начать отправте ваши контакты")
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Отправляя ваши контакты вы соглашаетесь с пользовательским соглашением \"Hellobox\":https://hellobox.uz/privacy-policy\n___*Сервис пока работает только на территории города Ташкента___")
			message.ReplyMarkup = reply
			msg.ParseMode = "markdown"
			bot.Send(message)
			bot.Send(msg)
			continue
		}
		user := database.FilterUser(models.User{ChatId: update.Message.From.ID})
		if user.Id == 0 || user.PartnerId == 0 {
			env.PartnerBot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "You are not authorized to use this bot"))
			continue
		}
		checkProduct(update.Message.Chat.ID, update.Message.From.ID, update.Message.Text)
		// sendPartnerMenu()
	}
}
