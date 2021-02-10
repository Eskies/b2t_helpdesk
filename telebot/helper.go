package telebot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func keyboardYaTidak() [][]tgbotapi.KeyboardButton {
	var ya tgbotapi.KeyboardButton
	ya.Text = "Iya"
	var tidak tgbotapi.KeyboardButton
	tidak.Text = "Tidak"

	var dmp [][]tgbotapi.KeyboardButton
	var dmp1 []tgbotapi.KeyboardButton
	dmp1 = append(dmp1, ya)
	dmp1 = append(dmp1, tidak)
	dmp = append(dmp, dmp1)
	return dmp

}
