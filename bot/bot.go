package bot

import (
	"fmt"
	"hellobox/database"
	"hellobox/env"
	"hellobox/models"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

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
			tgbotapi.KeyboardButton{Text: "🌟 Партнери"},
			tgbotapi.KeyboardButton{Text: "🔥 Топ продаж"},
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.KeyboardButton{Text: "📱 Контакты"},
			tgbotapi.KeyboardButton{Text: "🔑 Проверить продукт"},
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.KeyboardButton{Text: "🎁 Открыть полученный подарок"},
		),
	)
	message := tgbotapi.NewMessage(id, "Кто-то сегодня обрадуется 😍")
	message.ReplyMarkup = reply
	// del := tgbotapi.NewDeleteMessage(id, int(messageId))
	// env.Bot.Send(del)
	env.Bot.Send(message)
}

func hasContact(update tgbotapi.Update, fromPartner bool) bool {
	if update.Message == nil {
		return false
	}
	if update.Message.Contact != nil {
		user := models.User{
			Phone:     update.Message.Contact.PhoneNumber,
			Firstname: update.Message.Contact.FirstName,
			Lastname:  update.Message.Contact.LastName,
			TgId:      update.Message.From.ID,
			ChatId:    update.Message.Chat.ID,
		}
		if fromPartner {
			user.FromPartner = 1
		}
		database.CreateUser(user)
		// if err != nil {
		// 	env.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Вы "))
		// }
		if !fromPartner {
			sendMenu(update.Message.Chat.ID)
		} else {
			m := tgbotapi.NewMessage(update.Message.Chat.ID, "Now ask the admin to bind a partner to your account your id is: "+fmt.Sprint(update.Message.From.ID))
			m.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

			env.PartnerBot.Send(m)
		}
		return true
	}
	return false
}

func showCategories(id int64) {
	categories := database.GetCategories()
	markup := tgbotapi.NewInlineKeyboardMarkup()
	i := 1
	for i = 1; i < len(categories); i += 2 {
		count := len(database.GetProductsByCategory(categories[i-1].Id))
		count2 := len(database.GetProductsByCategory(categories[i].Id))
		markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s (%d)", categories[i-1].Name, count), fmt.Sprintf("%s#%d", "category", categories[i-1].Id)),
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s (%d)", categories[i].Name, count2), fmt.Sprintf("%s#%d", "category", categories[i].Id))))
	}
	if i <= len(categories) {
		markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(categories[i-1].Name, fmt.Sprintf("%s#%d", "partner", categories[i-1].Id))))
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
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s %d сум", el.Name, el.Price), fmt.Sprintf("%s#%d", "product", el.Id)),
		))
	}
	message := tgbotapi.NewMessage(chatId, "Выберите продукт")
	message.ReplyMarkup = markup
	del := tgbotapi.NewDeleteMessage(chatId, messageId)
	env.Bot.Send(del)
	env.Bot.Send(message)
}

func showPartner(chatId int64) {
	partner := database.GetPartner()
	markup := tgbotapi.NewInlineKeyboardMarkup()
	i := 1
	for i = 1; i < len(partner); i += 2 {
		markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(partner[i-1].Name, fmt.Sprintf("%s#%d", "partner", partner[i-1].Id)),
			tgbotapi.NewInlineKeyboardButtonData(partner[i].Name, fmt.Sprintf("%s#%d", "partner", partner[i].Id))))
	}
	if i <= len(partner) {
		markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(partner[i-1].Name, fmt.Sprintf("%s#%d", "partner", partner[i-1].Id))))
	}
	message := tgbotapi.NewMessage(chatId, "Выберите нужный продукт из списка партнеров")
	message.ReplyMarkup = markup
	env.Bot.Send(message)
}

func showSettings(chatId int64) {
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
		s := u.String()
		token := s[len(s)-12 : len(s)-1]
		realProduct := database.GetSingleProduct(productId)
		user.Cart.Products = append(user.Cart.Products, models.CartProduct{ProductId: realProduct.Id, Product: realProduct, CartId: user.CartId, Count: 1, Token: token, OptionIndex: 0})
	}
	database.EditUser(user)
	if isCart {
		showCart(user.ChatId, update.CallbackQuery.Message.MessageID, isCart, nil)
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
			//Remove product from list
			user.Cart.RemoveProduct(productId)
			database.DeleteCartProduct(product.Id)
			// database.ClearUserCart(user)
		} else {
			user.Cart.SetProduct(*product)
		}
		fmt.Println()
		fmt.Println()
		fmt.Println(len(user.Cart.Products))
		fmt.Println()
		fmt.Println()
	}
	database.EditUser(user)
	if isCart {
		showCart(user.ChatId, update.CallbackQuery.Message.MessageID, isCart, &user)
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
		env.Bot.Send(tgbotapi.NewCallbackWithAlert(update.CallbackQuery.ID, "Пожалуйста укажите количество, чтобы продолжить."))
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
			if user.Cart != nil && len(user.Cart.Products) > 0 {
				selectedOption = user.Cart.GetProduct(productId).OptionIndex
			}
		} else {
			if optionIndex == 1000 {
				pr := user.Cart.GetProduct(productId)
				pr.IsPresent = true
				pr.OptionIndex = 0
				database.EditCartProduct(*pr)
				fmt.Print("\n\n\nPRESENT\n\n\n")
				fmt.Printf("\n\n\n%s\n\n\n", strconv.FormatBool(pr.IsPresent))
				fmt.Println(pr)
			} else {
				selectedOption = uint(optionIndex)
				pr := user.Cart.GetProduct(productId)
				pr.OptionIndex = uint(optionIndex)
				pr.IsPresent = false
				if optionIndex == 1000 || pr.IsPresent {
					pr.OptionIndex = 0
				}
				database.EditCartProduct(*pr)
			}
		}
		isPresent := user.Cart.GetProduct(productId).IsPresent
		firstText := "1 ☑️"
		if selectedOption != 0 || optionIndex == 1000 || isPresent {
			firstText = "1"
		}
		row := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(firstText, fmt.Sprintf("option#0#%d", productId)))
		//Appending product options
		for i, _ := range product.Options {
			txt := fmt.Sprint(i + 2)
			if i+1 == int(selectedOption) && optionIndex != 1000 && !isPresent {
				txt = fmt.Sprintf("%d ☑️", i+2)
			}
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(txt, fmt.Sprintf("option#%d#%d", i+1, productId)))
		}

		markup.InlineKeyboard = append(markup.InlineKeyboard, row)
		p := user.Cart.GetProduct(productId)
		if p.IsPresent || optionIndex == 1000 || p.OptionIndex == 1000 {
			markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🎁 Отправить как подарок ☑️", fmt.Sprintf("select-present#%d", p.ProductId))))
		} else {
			markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🎁 Отправить как подарок", fmt.Sprintf("select-present#%d", p.ProductId))))
		}
		markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Оформить заказ", "show-cart"),
		))
	} else {
		markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🎁 Кастомизировать продукт", fmt.Sprintf("customize#%d", productId))))
	}
	markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Назад", fmt.Sprintf("product-back#%d", product.CategoryId)),
	))

	text := fmt.Sprintf("***%s***\nЦена:%d\n%s", product.Name, product.Price, product.Description)
	if isEdit {
		//!TODO COMPLETE
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
		if showOptions {
			url := product.ImageUrl
			if user.Cart != nil && len(user.Cart.Products) > 0 && selectedOption != 0 && optionIndex != 1000 {
				pr := user.Cart.GetProduct(productId)
				if pr.OptionIndex == 1000 {
					pi := database.GetPresentImage()
					url = pi.ImageUrl
				} else {
					url = pr.Product.Options[selectedOption-1].ImageUrl
				}
			}
			prod := user.Cart.GetProduct(productId)
			if prod.IsPresent || optionIndex == 1000 {
				pi := database.GetPresentImage()
				url = pi.ImageUrl
			}
			fmt.Println()
			fmt.Println()
			fmt.Println(selectedOption, prod.IsPresent, optionIndex)
			fmt.Println()
			fmt.Println()
			del := tgbotapi.NewDeleteMessage(user.ChatId, messageId)
			env.Bot.Send(del)
			file := tgbotapi.NewPhoto(chatId, tgbotapi.FileURL(url))
			file.ParseMode = "markdown"
			file.Caption = text
			file.ReplyMarkup = markup
			env.Bot.Send(file)
			return
		}
		edit := tgbotapi.EditMessageCaptionConfig(tgbotapi.EditMessageCaptionConfig{BaseEdit: tgbotapi.BaseEdit{ChatID: user.ChatId, MessageID: messageId}})
		edit.ReplyMarkup = &markup
		edit.ParseMode = "markdown"
		edit.Caption = text
		if !isEdit {
			del := tgbotapi.NewDeleteMessage(chatId, messageId)
			env.Bot.Send(del)
		}
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
		if optionIndex == 1000 || pr.IsPresent {
			pr.OptionIndex = 0
		}
		database.EditCartProduct(*pr)
	}
	url := product.ImageUrl
	if user.Cart != nil && len(user.Cart.Products) > 0 && selectedOption != 0 {
		if pr := user.Cart.GetProduct(productId); !pr.IsPresent && pr.OptionIndex != 1000 {
			url = pr.Product.Options[selectedOption-1].ImageUrl
		} else {
			pi := database.GetPresentImage()
			url = pi.ImageUrl
		}
	}
	file := tgbotapi.NewPhoto(chatId, tgbotapi.FileURL(url))
	file.ParseMode = "markdown"
	file.Caption = text
	file.ReplyMarkup = markup
	del := tgbotapi.NewDeleteMessage(chatId, messageId)
	env.Bot.Send(del)
	env.Bot.Send(file)
}

func showCart(chatId int64, messageId int, isEdit bool, usr *models.User) {
	user := models.User{TgId: chatId}
	if usr != nil {
		user = *usr
	} else {
		user = database.FilterUser(user)
	}

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
		txt = fmt.Sprintf("%s\n***%s***\n└  %s  %d x %d = %d", txt, el.Product.Name, el.Product.Name, el.Count, el.Product.Price, el.Product.Price*el.Count)
	}
	percent := database.GetProfitPercent()
	markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("✅ Оформить заказ", "order")))
	total := float32(user.Cart.CartTotal())
	commission := float32(user.Cart.CartTotal()) * float32(float32(percent.Percent)/100.0)
	txt = fmt.Sprintf("%s\n\nКомиссия: %d%%\n\nВсего: %0.0f сум", txt, percent.Percent, total+commission)
	if isEdit {
		msg := tgbotapi.NewEditMessageTextAndMarkup(chatId, messageId, txt, markup)
		msg.ParseMode = "markdown"
		env.Bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(chatId, txt)
		msg.ReplyMarkup = markup
		msg.ParseMode = "markdown"
		env.Bot.Send(msg)
		del := tgbotapi.NewDeleteMessage(chatId, messageId)
		env.Bot.Send(del)
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
		markup := tgbotapi.InlineKeyboardMarkup{}
		if news.ProductId != 0 {
			markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(news.Product.Name, fmt.Sprintf("%s#%d", "product", news.ProductId))))
		}
		if news.PartnerId != 0 && news.ProductId == 0 {
			markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(news.Partner.Name, fmt.Sprintf("%s#%d", "partner", news.PartnerId))))
		}
		if len(markup.InlineKeyboard) > 0 {
			file.ReplyMarkup = markup
		}
		env.Bot.Send(file)
	}
}

func makeOrder(update tgbotapi.Update) {
	user := database.FilterUser(models.User{TgId: update.CallbackQuery.From.ID})
	database.CreateOrder(models.Order{
		Cart:   *user.Cart,
		UserId: user.Id,
	})
	items := []tgbotapi.LabeledPrice{}
	txt := "Оплата"
	for _, el := range user.Cart.Products {
		items = append(items, tgbotapi.LabeledPrice{Label: el.Product.Name, Amount: int(el.Product.Price) * 100 * int(el.Count)})
	}
	percent := database.GetProfitPercent()
	total := float32(user.Cart.CartTotal())
	commission := float32(user.Cart.CartTotal()) * float32(float32(percent.Percent))
	println(commission)
	items = append(items, tgbotapi.LabeledPrice{Label: fmt.Sprintf("Комиссия: %d%%", percent.Percent), Amount: int(commission)})
	token := "387026696:LIVE:61d30e670f5ef6a30739d8c3"
	// token := "371317599:TEST:1638986618188"
	in := tgbotapi.NewInvoice(update.CallbackQuery.Message.Chat.ID, "Hellobox", txt, user.Cart.Products[0].Token, token, user.Cart.Products[0].Token, "UZS", items)
	in.SuggestedTipAmounts = []int{}
	pay := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("💳 %0.0f", total+(commission/100.0)), "")
	pay.Pay = true
	in.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(pay))
	del := tgbotapi.NewDeleteMessage(update.CallbackQuery.From.ID, int(update.CallbackQuery.Message.MessageID))
	env.Bot.Send(del)
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
	if update.Message.Text == "" {
		return
	}
	product := database.GetCartProductsByToken(update.Message.Text)
	if product.Id == 0 {
		env.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Недействителен"))
		return
	}
	file := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FileURL(product.Product.ImageUrl))
	file.ParseMode = "markdown"
	text := fmt.Sprintf("***%s***\n%s", product.Product.Name, product.Product.Description)
	file.Caption = text
	env.Bot.Send(file)
}

func showTopProducts(update tgbotapi.Update) {
	products := database.GetTopProducts()
	markup := tgbotapi.NewInlineKeyboardMarkup()
	for _, el := range products {
		markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s %d сум", el.Name, el.Price), fmt.Sprintf("%s#%d", "product", el.Id)),
		))
	}
	message := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите продукт")
	message.ReplyMarkup = markup
	env.Bot.Send(message)
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
			//Sending prompt
			txt := "Ваш заказ выполнен ✅\nМожете переслать изображение продукта с кодом любому человеку 😊"
			prompt := tgbotapi.NewMessage(update.PreCheckoutQuery.From.ID, txt)
			env.Bot.Send(prompt)

			//Sending tokens
			for _, el := range user.Cart.Products {
				//Product data
				selectedOption := el.OptionIndex
				url := el.Product.ImageUrl
				if selectedOption != 0 {
					if selectedOption == 1000 {
						pr := database.GetPresentImage()
						url = pr.ImageUrl
					} else {
						url = el.Product.Options[selectedOption-1].ImageUrl
					}
				}
				//Sending photo
				file := tgbotapi.NewPhoto(user.TgId, tgbotapi.FileURL(url))
				file.ParseMode = "MarkdownV2"
				text := fmt.Sprintf("***%s***\n%s", el.Product.Name, el.Product.Description)
				if el.OptionIndex == 1000 {
					text = fmt.Sprintf("***Вам отправили подарок 🎁***\nЧтобы его открыть переходите  по ссылке в Телеграм бот «Hellobox \\(https://t\\.me/helloboxbot\\)»\\.И напишите «Код активации» в раздел: \n***«🎁 Открыть полученный подарок»\\.***")
				}

				t := time.Now()
				from := t.Format("02\\.01\\.2006")
				t = t.AddDate(0, 0, int(el.Product.ExpiresIn))
				to := t.Format("02\\.01\\.2006")

				txt := fmt.Sprintf("⏰Период активации: %s\\-%s\n 🐼 Количество:%d\n🔑Код активации:  ||***%s***||\n_\\*Не показывайте код другим людям, его можете активировать только вы или знающий его человек;\n\\*Сервис пока работает только на территории города Ташкента_", from, to, el.Count, el.Token)
				file.Caption = text + "\n" + txt
				env.Bot.Send(file)
				// msg := tgbotapi.NewMessage(update.PreCheckoutQuery.From.ID, )
				// msg.ParseMode = "MarkdownV2"
				// env.Bot.Send(msg)
			}
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
			case "show-cart":
				showCart(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, false, nil)
			case "select-present":
				productId, err := strconv.ParseUint(s[1], 10, 32)
				if err != nil {
					continue
				}
				showProductDetails(update, uint(productId), update.CallbackQuery.Message.MessageID, true, nil, true, 1000)
			}

			continue
		}
		if update.Message == nil {
			continue
		}

		if hasContact(update, false) {
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
		case "🛍 Каталог":
			showCategories(update.Message.Chat.ID)
		case "🌟 Партнери":
			showPartner(update.Message.Chat.ID)
		case "☎️ Обратная связь":
			showSettings(update.Message.Chat.ID)
		case "🛒 Корзина":
			showCart(update.Message.Chat.ID, -1, false, nil)
		case "🎁 Открыть полученный подарок":
			openRecievedGift(update)
		case "🔑 Проверить продукт":
			// orderHistory()
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пожалуйста напишите код активации, чтобы проверить статус продукта")
			env.Bot.Send(msg)
		case "📱 Контакты":
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "По всем вопросам обращаться к оператору @Helloboxuz"))
		case "🔥 Топ продаж":
			showTopProducts(update)
		default:
			showProductDetailsByToken(update)
		}

	}
}
