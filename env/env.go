package env

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gorilla/mux"
)

var (
	Bot        *tgbotapi.BotAPI
	PartnerBot *tgbotapi.BotAPI
	Router     *mux.Router
)
