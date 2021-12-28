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
			tgbotapi.KeyboardButton{Text: "📜 История заказов"},
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.KeyboardButton{Text: "🎁 Открыть полученный подарок"},
		),
	)
	message := tgbotapi.NewMessage(id, "Кто то сегодня обрадуется")
	message.ReplyMarkup = reply
	env.Bot.Send(message)
}

func hasContact(update tgbotapi.Update) bool {
	if update.Message.Contact != nil {
		database.CreateUser(models.User{
			Phone:     update.Message.Contact.PhoneNumber,
			Firstname: update.Message.Contact.FirstName,
			Lastname:  update.Message.Contact.LastName,
			TgId:      update.Message.From.ID,
			ChatId:    update.Message.Chat.ID,
		})
		// if err != nil {
		// 	env.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Вы "))
		// }
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
	message := tgbotapi.NewMessage(id, "Выберите нужную категорию ☺️ ")
	message.ReplyMarkup = markup
	env.Bot.Send(message)

}

func showProducts(chatId int64, categoryId uint, messageId int, partnerId uint) {
	var products []models.Product
	if partnerId == 0 {
		products = database.GetProductsByCategory(categoryId)
	} else {
		products = database.GetProductsByPartner(partnerId)
	}
	markup := tgbotapi.NewInlineKeyboardMarkup()
	for _, el := range products {
		markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(el.Name, fmt.Sprintf("%s#%d", "product", el.Id)),
		))
	}
	message := tgbotapi.NewMessage(chatId, "Выберите продукт")
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
	message := tgbotapi.NewMessage(chatId, "Выберите нужный продукт из списка партнеров")
	message.ReplyMarkup = markup
	env.Bot.Send(message)
}

func showSettings(chatId int64, messageId int) {
	settings := database.GetSettings()
	for i := 1; i < len(settings); i += 2 {
	}
	message := tgbotapi.NewMessage(chatId, "По всем вопросам обращаться к оператору @Helloboxuz")
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
			user.Cart = &models.Cart{}
		}
		u, _ := uuid.NewV4()
		realProduct := database.GetSingleProduct(productId)
		user.Cart.Products = append(user.Cart.Products, models.CartProduct{ProductId: realProduct.Id, Product: realProduct, CartId: user.CartId, Count: 1, Token: u.String(), OptionIndex: 0})
	}
	database.EditUser(user)
	if isCart {
		showCart(user.ChatId, update.CallbackQuery.Message.MessageID, isCart)
		return
	}
	showProductDetails(update, productId, update.CallbackQuery.Message.MessageID, true, &user, false, -1)
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
	showProductDetails(update, productId, update.CallbackQuery.Message.MessageID, true, &user, false, -1)
}

func showProductDetails(update tgbotapi.Update, productId uint, messageId int, isEdit bool, user *models.User, showOptions bool, optionIndex int) {
	product := database.GetSingleProduct(productId)
	chatId := update.CallbackQuery.Message.Chat.ID
	if user == nil {
		user = &models.User{ChatId: chatId}
		*user = database.FilterUser(*user)
	}
	if showOptions && (user.Cart == nil || len(user.Cart.Products) <= 0) {
		env.Bot.Send(tgbotapi.NewCallbackWithAlert(update.CallbackQuery.ID, "Пожалуйста, сначала выберите хотя бы один продукт"))
		return
	}
	count := 0
	if user.Cart != nil && len(user.Cart.Products) > 0 {
		p := user.Cart.GetProduct(productId)
		if p != nil {
			count = int(p.Count)
		}
	}
	var markup tgbotapi.InlineKeyboardMarkup
	if !showOptions {
		markup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("🛒%d", user.Cart.CartTotal()), "cart"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("➖", fmt.Sprintf("minus#%d", productId)),
				tgbotapi.NewInlineKeyboardButtonData(fmt.Sprint(count), "none"),
				tgbotapi.NewInlineKeyboardButtonData("➕", fmt.Sprintf("plus#%d", productId)),
			),
		)
	}
	if showOptions {
		var selectedOption uint
		if optionIndex == -1 {
			selectedOption = user.Cart.GetProduct(productId).OptionIndex
		} else {
			selectedOption = uint(optionIndex)
			pr := user.Cart.GetProduct(productId)
			pr.OptionIndex = uint(optionIndex)
			database.EditCartProduct(*pr)
		}
		firstText := "1 ✅"
		if selectedOption != 0 {
			firstText = "1"
		}
		row := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(firstText, fmt.Sprintf("option#0#%d", productId)))
		//Appending product options
		for i, _ := range product.Options {
			txt := fmt.Sprint(i + 2)
			if i+1 == int(selectedOption) {
				txt = fmt.Sprintf("%d ✅", i+2)
			}
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(txt, fmt.Sprintf("option#%d#%d", i+1, productId)))
		}

		markup.InlineKeyboard = append(markup.InlineKeyboard, row)
		markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎁 Отправить как подарок", "present")))
	} else {
		markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🎁 Кастомизировать продукт", fmt.Sprintf("customize#%d", productId))))
	}
	markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Назад", fmt.Sprintf("product-back#%d", product.CategoryId)),
	))

	text := fmt.Sprintf("***%s***\nЦена:%d\n%s", product.Name, product.Price, product.Description)
	if isEdit {
		//!TODO COMPLETE
		// var selectedOption uint
		// if showOptions {
		// 	url := product.ImageUrl
		// 	if user.Cart != nil && len(user.Cart.Products) > 0 && selectedOption != 0 {
		// 		pr := user.Cart.GetProduct(productId)
		// 		url = pr.Product.Options[selectedOption-1].ImageUrl
		// 	}
		// 	edit := tgbotapi.EditMessageMediaConfig(tgbotapi.EditMessageMediaConfig{BaseEdit: tgbotapi.BaseEdit{ChatID: user.ChatId, MessageID: messageId}, Media: tgbotapi.FileURL(url)})
		// 	edit.ReplyMarkup = &markup
		// 	env.Bot.Send(edit)
		// 	return
		// }
		edit := tgbotapi.EditMessageCaptionConfig(tgbotapi.EditMessageCaptionConfig{BaseEdit: tgbotapi.BaseEdit{ChatID: user.ChatId, MessageID: messageId}})
		edit.ReplyMarkup = &markup
		edit.ParseMode = "markdown"
		edit.Caption = text
		env.Bot.Send(edit)
		return
	}

	var selectedOption uint
	if optionIndex == -1 {
		if user.Cart == nil || len(user.Cart.Products) <= 0 {
			selectedOption = 0
		} else {
			p := user.Cart.GetProduct(productId)
			if p != nil {
				selectedOption = p.OptionIndex
			} else {
				selectedOption = 0
			}
		}
	} else {
		selectedOption = uint(optionIndex)
		pr := user.Cart.GetProduct(productId)
		pr.OptionIndex = uint(optionIndex)
		database.EditCartProduct(*pr)
	}
	url := product.ImageUrl
	if user.Cart != nil && len(user.Cart.Products) > 0 && selectedOption != 0 {
		pr := user.Cart.GetProduct(productId)
		url = pr.Product.Options[selectedOption-1].ImageUrl
	}
	file := tgbotapi.NewPhoto(chatId, tgbotapi.FileURL(url))
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
	items := []tgbotapi.LabeledPrice{}
	for _, el := range user.Cart.Products {
		items = append(items, tgbotapi.LabeledPrice{Label: el.Product.Name, Amount: int(el.Product.Price) * 100 * int(el.Count)})
	}
	in := tgbotapi.NewInvoice(update.CallbackQuery.Message.Chat.ID, "PAyment", "Pay plz", user.Cart.Products[0].Token, "371317599:TEST:1638986618188", user.Cart.Products[0].Token, "UZS", items)
	in.SuggestedTipAmounts = []int{}
	env.Bot.Send(in)
}

func sendAsPresent(update tgbotapi.Update) {
	// msg := tgbotapi.NewMessage(chat,"Перешлите это сообщение человеку, которого вы хотите представить")
}

func openRecievedGift(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пожалуйста напишите полученный код")
	env.Bot.Send(msg)
}

func showProductDetailsByToken(update tgbotapi.Update) {
	product := database.GetCartProductsByToken(update.Message.Text)
	file := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FileURL(product.Product.ImageUrl))
	file.ParseMode = "markdown"
	text := fmt.Sprintf("***%s***\n%s", product.Product.Name, product.Product.Description)
	file.Caption = text
	env.Bot.Send(file)
}

func orderHistory(update tgbotapi.Update) {
	// user := database.FilterUser(models.User{TgId: update.Message.From.ID})
	// history := database.GetOrderHistory(user.Id)
}

func HandleBot() {
	bot, err := tgbotapi.NewBotAPI("2105220707:AAGzYOqoUOEwXgupuKseRbVQGG5Y5Tqluv8")
	env.Bot = bot
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)

	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.PreCheckoutQuery != nil {
			user := database.FilterUser(models.User{TgId: update.PreCheckoutQuery.From.ID})
			if user.Cart == nil {
				bot.Send(tgbotapi.PreCheckoutConfig{PreCheckoutQueryID: update.PreCheckoutQuery.ID, OK: false, ErrorMessage: "Something went wrong"})
				continue
			}
			bot.Send(tgbotapi.PreCheckoutConfig{PreCheckoutQueryID: update.PreCheckoutQuery.ID, OK: true})
			product := database.GetCartProductsByToken(update.PreCheckoutQuery.InvoicePayload)
			selectedOption := product.OptionIndex
			url := product.Product.ImageUrl
			if selectedOption != 0 {
				url = product.Product.Options[selectedOption-1].ImageUrl
			}
			file := tgbotapi.NewPhoto(user.ChatId, tgbotapi.FileURL(url))
			file.ParseMode = "markdown"
			text := fmt.Sprintf("***%s***\n%s", product.Product.Name, product.Product.Description)
			file.Caption = text
			env.Bot.Send(file)
			txt := "Ваш заказ выполнен."
			msg := tgbotapi.NewMessage(update.PreCheckoutQuery.From.ID, fmt.Sprintf("%s\nВаш токен: ***%s***", txt, user.Cart.Products[0].Token))
			msg.ParseMode = "markdown"
			env.Bot.Send(msg)
			database.ClearUserCart(user)
			continue
		}
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
				showProducts(update.CallbackQuery.Message.Chat.ID, uint(categoryId), update.CallbackQuery.Message.MessageID, 0)
			case "product":
				//Product details
				productId, err := strconv.ParseUint(s[1], 10, 32)
				if err != nil {
					continue
				}
				//Show product details
				showProductDetails(update, uint(productId), update.CallbackQuery.Message.MessageID, false, nil, false, -1)
			case "product-back":
				categoryId, err := strconv.ParseUint(s[1], 10, 32)
				if err != nil {
					continue
				}
				showProducts(update.CallbackQuery.Message.Chat.ID, uint(categoryId), update.CallbackQuery.Message.MessageID, 0)
			case "order":
				makeOrder(update)
			case "partner":
				partnerId, err := strconv.ParseUint(s[1], 10, 32)
				if err != nil {
					continue
				}
				showProducts(update.CallbackQuery.Message.Chat.ID, 0, update.CallbackQuery.Message.MessageID, uint(partnerId))
			case "customize":
				productId, err := strconv.ParseUint(s[1], 10, 32)
				if err != nil {
					continue
				}
				showProductDetails(update, uint(productId), update.CallbackQuery.Message.MessageID, true, nil, true, -1)
			case "option":
				optionIndex, err := strconv.ParseUint(s[1], 10, 32)
				if err != nil {
					continue
				}
				productId, err := strconv.ParseUint(s[2], 10, 32)
				if err != nil {
					continue
				}
				showProductDetails(update, uint(productId), update.CallbackQuery.Message.MessageID, true, nil, true, int(optionIndex))
			case "present":
				sendAsPresent(update)
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
					tgbotapi.NewKeyboardButtonContact("Отправить ваши контакы"),
				),
			)
			message := tgbotapi.NewMessage(update.Message.Chat.ID, "Чтобы начать отправте ваши контакты")
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Отправляя ваши контакты вы соглашаетесь политикой конфиденциальности Hellobox\nhttps://hellobox.uz/privacy-policy")
			message.ReplyMarkup = reply
			bot.Send(message)
			bot.Send(msg)
		case "🛍 Каталог":
			showCategories(update.Message.Chat.ID)
		case "👥 Партнери":
			showPartner(update.Message.Chat.ID, update.Message.MessageID)
		case "☎️ Обратная связь":
			showSettings(update.Message.Chat.ID, update.Message.MessageID)
		case "🛒 Корзина":
			showCart(update.Message.Chat.ID, -1, false)
		case "🎁 Открыть полученный подарок":
			openRecievedGift(update)
		case "📜 История заказов":
			// orderHistory()
		}
		showProductDetailsByToken(update)
	}
}
