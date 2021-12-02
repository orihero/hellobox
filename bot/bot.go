package bot

import (
	"fmt"
	"hellobox/database"
	"hellobox/env"
	"hellobox/models"
	"io/ioutil"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	uuid "github.com/nu7hatch/gouuid"
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
			TgId:      update.Message.From.ID,
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

func showCategories(id int64) {
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
	message := tgbotapi.NewMessage(chatId, "Buyurtmani tanlang!")
	message.ReplyMarkup = markup
	env.Bot.Send(message)
}

func showPartner(chatId int64, messageId int) {
	partner := database.GetPartner()
	markup := tgbotapi.NewInlineKeyboardMarkup()
	for i := 1; i < len(partner); i += 2 {
		markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(partner[i-1].Name, fmt.Sprintf("%s#%d", "partner", partner[i-1].Id)),
			tgbotapi.NewInlineKeyboardButtonData(partner[i].Name, fmt.Sprintf("%s#%d", "partner", partner[i].Id))))
	}
	message := tgbotapi.NewMessage(chatId, "Biz bilan hamkorlar")
	message.ReplyMarkup = markup
	env.Bot.Send(message)
}

func showSettings(chatId int64, messageId int) {
	settings := database.GetSettings()
	for i := 1; i < len(settings); i += 2 {
	}
	message := tgbotapi.NewMessage(chatId, "Biz haqimizda http://www.telegram.org/helloboxuz")
	env.Bot.Send(message)
}

func incrementProduct(update tgbotapi.Update, productId uint, isCart bool) {
	user := models.User{TgId: update.CallbackQuery.From.ID}
	user = database.FilterUser(user)
	product := user.Cart.GetProduct(productId)
	if product != nil {
		product.Count++
		user.Cart.SetProduct(*product)
	} else {
		if user.Cart == nil || len(user.Cart.Products) == 0 {
			u, _ := uuid.NewV4()
			user.Cart = &models.Cart{Token: u.String()}
		}
		realProduct := database.GetSingleProduct(productId)
		user.Cart.Products = append(user.Cart.Products, models.CartProduct{ProductId: realProduct.Id, Product: realProduct, CartId: user.CartId, Count: 1})
	}
	database.EditUser(user)
	if isCart {
		showCart(user.ChatId, update.CallbackQuery.Message.MessageID, isCart)
		return
	}
	showProductDetails(user.ChatId, productId, update.CallbackQuery.Message.MessageID, true, &user)
}

func decrementProduct(update tgbotapi.Update, productId uint, isCart bool) {
	user := models.User{TgId: update.CallbackQuery.From.ID}
	user = database.FilterUser(user)
	product := user.Cart.GetProduct(productId)
	if product != nil && product.Count > 0 {
		product.Count--
		if product.Count == 0 {
			// Remove product from list
			// user.Cart.RemoveProduct(productId)
			database.ClearUserCart(user)
		}
		user.Cart.SetProduct(*product)
	}
	database.EditUser(user)
	if isCart {
		showCart(user.ChatId, update.CallbackQuery.Message.MessageID, isCart)
		return
	}
	showProductDetails(user.ChatId, productId, update.CallbackQuery.Message.MessageID, true, &user)
}

func showProductDetails(chatId int64, productId uint, messageId int, isEdit bool, user *models.User) {
	product := database.GetSingleProduct(productId)
	if user == nil {
		user = &models.User{ChatId: chatId}
		*user = database.FilterUser(*user)
	}
	count := 0
	if user.Cart != nil && len(user.Cart.Products) > 0 {
		count = int(user.Cart.GetProduct(productId).Count)
	}
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("🛒%d", user.Cart.CartTotal()), "cart"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➖", fmt.Sprintf("minus#%d", productId)),
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprint(count), "none"),
			tgbotapi.NewInlineKeyboardButtonData("➕", fmt.Sprintf("plus#%d", productId)),
		),

		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", fmt.Sprintf("product-back#%d", product.CategoryId)),
		),
	)
	text := fmt.Sprintf("***%s***\nЦена:%d\n%s", product.Name, product.Price, product.Description)
	if isEdit {
		edit := tgbotapi.EditMessageCaptionConfig(tgbotapi.EditMessageCaptionConfig{BaseEdit: tgbotapi.BaseEdit{ChatID: user.ChatId, MessageID: messageId}})
		edit.ReplyMarkup = &markup
		edit.ParseMode = "markdown"
		edit.Caption = text
		env.Bot.Send(edit)
		return
	}
	file := tgbotapi.NewPhoto(chatId, tgbotapi.FileURL(product.ImageUrl))
	file.ParseMode = "markdown"
	file.Caption = text
	file.ReplyMarkup = markup
	env.Bot.Send(file)
}

func showCart(chatId int64, messageId int, isEdit bool) {
	user := models.User{TgId: chatId}
	user = database.FilterUser(user)

	txt := "***Корзинка***"
	if user.Cart == nil || len(user.Cart.Products) <= 0 {
		msg := tgbotapi.NewMessage(chatId, fmt.Sprintf("%s\nВаша корзинка пустая", txt))
		msg.ParseMode = "markdown"
		env.Bot.Send(msg)
		if isEdit {
			env.Bot.Send(tgbotapi.NewDeleteMessage(chatId, messageId))
		}
		return
	}
	markup := tgbotapi.NewInlineKeyboardMarkup()
	for _, el := range user.Cart.Products {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➖", fmt.Sprintf("cart-minus#%d", el.ProductId)),
			tgbotapi.NewInlineKeyboardButtonData(el.Product.Name, "none"),
			tgbotapi.NewInlineKeyboardButtonData("➕", fmt.Sprintf("cart-plus#%d", el.ProductId)),
		)
		markup.InlineKeyboard = append(markup.InlineKeyboard, row)
		markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Заказать", "order")))
		txt = fmt.Sprintf("%s\n***%s***\n└  %s  %d x %d = %d", txt, el.Product.Name, el.Product.Name, el.Count, el.Product.Price, el.Product.Price*el.Count)
	}
	txt = fmt.Sprintf("%s\n\nUmumiy:%d so'm", txt, user.Cart.CartTotal())
	if isEdit {
		msg := tgbotapi.NewEditMessageTextAndMarkup(chatId, messageId, txt, markup)
		msg.ParseMode = "markdown"
		env.Bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(chatId, txt)
		msg.ReplyMarkup = markup
		msg.ParseMode = "markdown"
		env.Bot.Send(msg)
	}
}

func SendNews(news models.News) {
	users := database.GetUsers()
	for _, el := range users {
		//TODO FIX
		str := strings.Split(news.ImageUrl, "/")
		last := str[len(str)-1]
		body, _ := ioutil.ReadFile(fmt.Sprintf("./public/uploads/%s", last))
		file := tgbotapi.NewPhoto(el.ChatId, tgbotapi.FileBytes{Bytes: body})
		file.Caption = news.Description
		env.Bot.Send(file)
	}
}

func makeOrder(update tgbotapi.Update) {
	user := database.FilterUser(models.User{TgId: update.CallbackQuery.From.ID})
	database.CreateOrder(models.Order{
		Cart: *user.Cart,
	})
	txt := "Ваш заказ выполнен."
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("%s\nВаш токен: ***%s***", txt, user.Cart.Token))
	msg.ParseMode = "markdown"
	env.Bot.Send(msg)
	database.ClearUserCart(user)
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
			case "minus":
				productId, err := strconv.ParseUint(s[1], 10, 32)
				if err != nil {
					continue
				}
				decrementProduct(update, uint(productId), false)
			case "plus":
				productId, err := strconv.ParseUint(s[1], 10, 32)
				if err != nil {
					continue
				}
				incrementProduct(update, uint(productId), false)
			case "cart-minus":
				productId, err := strconv.ParseUint(s[1], 10, 32)
				if err != nil {
					continue
				}
				decrementProduct(update, uint(productId), true)
			case "cart-plus":
				productId, err := strconv.ParseUint(s[1], 10, 32)
				if err != nil {
					continue
				}
				incrementProduct(update, uint(productId), true)
			case "category":
				//Show products
				categoryId, err := strconv.ParseUint(s[1], 10, 32)
				if err != nil {
					continue
				}
				showProducts(update.CallbackQuery.Message.Chat.ID, uint(categoryId), update.CallbackQuery.Message.MessageID)
			case "product":
				//Product details
				productId, err := strconv.ParseUint(s[1], 10, 32)
				if err != nil {
					continue
				}
				//Show product details
				showProductDetails(update.CallbackQuery.Message.Chat.ID, uint(productId), update.CallbackQuery.Message.MessageID, false, nil)
			case "product-back":
				categoryId, err := strconv.ParseUint(s[1], 10, 32)
				if err != nil {
					continue
				}
				showProducts(update.CallbackQuery.Message.Chat.ID, uint(categoryId), update.CallbackQuery.Message.MessageID)
			case "order":
				makeOrder(update)
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
			showCategories(update.Message.Chat.ID)
		case "👥 Партнери":
			showPartner(update.Message.Chat.ID, update.Message.MessageID)
		case "☎️ Обратная связь":
			showSettings(update.Message.Chat.ID, update.Message.MessageID)
		case "🛒 Корзина":
			showCart(update.Message.Chat.ID, -1, false)
		}

	}
}
