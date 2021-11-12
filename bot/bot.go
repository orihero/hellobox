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

func sendMenu(id int64) {
	reply := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.KeyboardButton{Text: "🛍 Каталог"},
			tgbotapi.KeyboardButton{Text: "🛒 Корзина"},
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.KeyboardButton{Text: "👥 Партнери"},
			tgbotapi.KeyboardButton{Text: "🔥 Топ продаж"},
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.KeyboardButton{Text: "☎️ Обратная связь"},
		),
	)
	message := tgbotapi.NewMessage(id, "Kerakli funksiyani tanlang")
	message.ReplyMarkup = reply
	env.Bot.Send(message)
}

func hasContact(update tgbotapi.Update) bool {
	if update.Message.Contact != nil {
		err := database.CreateUser(models.User{
			Phone:     update.Message.Contact.PhoneNumber,
			Firstname: update.Message.Contact.FirstName,
			Lastname:  update.Message.Contact.LastName,
			UserId:    update.Message.From.ID,
			ChatId:    update.Message.Chat.ID,
		})
		if err != nil {
			env.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "User already exist"))
		}
		sendMenu(update.Message.Chat.ID)
		return true
	}
	return false
}

func showList(id int64) {
	categories := database.GetCategories()
	markup := tgbotapi.NewInlineKeyboardMarkup()
	for i := 1; i < len(categories); i += 2 {
		markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(categories[i-1].Name, fmt.Sprintf("%s#%d", "category", categories[i-1].Id)),
			tgbotapi.NewInlineKeyboardButtonData(categories[i].Name, fmt.Sprintf("%s#%d", "category", categories[i].Id))))
	}
	message := tgbotapi.NewMessage(id, "Menyuga hush kelibsiz ☺️ ")
	message.ReplyMarkup = markup
	env.Bot.Send(message)

}

func showProducts(chatId int64, categoryId uint, messageId int) {
	products := database.GetProductsByCategory(categoryId)
	markup := tgbotapi.NewInlineKeyboardMarkup()
	for _, el := range products {
		markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(el.Name, fmt.Sprintf("%s#%d", "product", el.Id)),
		))
	}
	message := tgbotapi.NewEditMessageTextAndMarkup(chatId, messageId, "Buyurtmani tanlang!", markup)
	env.Bot.Send(message)
}

func showProductDetails(chatId int64, productId uint, messageId int) {
	product := database.GetSingleProduct(productId)
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("🛒%d", product.Price), "cart"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("-", "like"),
			tgbotapi.NewInlineKeyboardButtonData("1", "like"),
			tgbotapi.NewInlineKeyboardButtonData("+", "like"),
		),

		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", fmt.Sprintf("product-back#%d", product.CategoryId)),
		),
	)
	text := fmt.Sprintf("***%s***\n```Цена:%d```\n%s", product.Name, product.Price, product.Description)
	file := tgbotapi.NewPhoto(chatId, tgbotapi.FileURL(product.ImageUrl))
	file.ParseMode = "markdown"
	file.Caption = text
	file.ReplyMarkup = markup
	env.Bot.Send(file)
}

func SendNews(news models.News) {
	users := database.GetUsers()
	for _, el := range users {
		file := tgbotapi.NewPhoto(el.ChatId, tgbotapi.FileURL(news.ImageUrl))
		file.Caption = news.Description
		env.Bot.Send(file)
	}
}

func HandleBot() {
	bot, err := tgbotapi.NewBotAPI("2105220707:AAEdCoI3Zbd0MY0UHYnjB5qeesHBVOSN1UU")
	env.Bot = bot
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)

	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.CallbackQuery != nil {
			println(update.CallbackQuery.Data)
			s := strings.Split(update.CallbackQuery.Data, "#")
			switch s[0] {
			case "category":
				//Show products
				// println(s[1])
				categoryId, err := strconv.ParseUint(s[1], 10, 32)
				if err != nil {
					continue
				}
				showProducts(update.CallbackQuery.Message.Chat.ID, uint(categoryId), update.CallbackQuery.Message.MessageID)
			case "product":
				productId, err := strconv.ParseUint(s[1], 10, 32)
				if err != nil {
					continue
				}
				//Show product details
				showProductDetails(update.CallbackQuery.Message.Chat.ID, uint(productId), update.CallbackQuery.Message.MessageID)
			case "product-back":
				categoryId, err := strconv.ParseUint(s[1], 10, 32)
				if err != nil {
					continue
				}
				showProducts(update.CallbackQuery.Message.Chat.ID, uint(categoryId), update.CallbackQuery.Message.MessageID)
			}
			continue
		}
		if update.Message == nil {
			continue
		}

		if hasContact(update) {
			continue
		}

		switch update.Message.Text {
		case "/start":
			reply := tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButtonContact("Send your contacts"),
				),
			)
			message := tgbotapi.NewMessage(update.Message.Chat.ID, "Please share your contacts")
			message.ReplyMarkup = reply
			bot.Send(message)
		case "🛍 Каталог":
			showList(update.Message.Chat.ID)
		}
	}
}
