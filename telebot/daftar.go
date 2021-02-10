package telebot

import (
	"b2t_helpdesk/injector"
	"encoding/json"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/tidwall/gjson"
)

func cmdDaftar(msg *tgbotapi.Message, di *injector.Injector, init bool) {
	if init {
		di.Enqueue(tgbotapi.NewMessage(msg.Chat.ID, di.Config.Get("pesan").Get("daftar").Get("pembuka").String()))
		di.Enqueue(tgbotapi.NewMessage(msg.Chat.ID, di.Config.Get("pesan").Get("daftar").Get("noregistrasi").String()))

		var cmdaction injector.CmdAction
		cmdaction.Cmd = "daftar"
		cmdaction.StepMax = 3
		cmdaction.StepNow = 0
		cmdaction.DataPerStep = append(cmdaction.DataPerStep, "start")

		json, _ := json.Marshal(cmdaction)

		di.SetRedisCmd(msg.From.ID, string(json), 5*time.Minute)
	} else {
		data := gjson.Parse(di.CmdData(msg.From.ID))
		switch data.Get("stepnow").Int() {
		case 0:
			var cmdaction injector.CmdAction
			cmdaction.Cmd = "daftar"
			cmdaction.StepMax = 3
			cmdaction.StepNow = 1
			data.Get("data").ForEach(func(key, value gjson.Result) bool {
				cmdaction.DataPerStep = append(cmdaction.DataPerStep, value.String())
				return true
			})
			cmdaction.DataPerStep = append(cmdaction.DataPerStep, msg.Text)
			di.DelRedisCmd(msg.From.ID)
			json, _ := json.Marshal(cmdaction)
			di.SetRedisCmd(msg.From.ID, string(json), 5*time.Minute)

			di.Enqueue(tgbotapi.NewMessage(msg.Chat.ID, di.Config.Get("pesan").Get("daftar").Get("nama").String()))
		case 1:
			var cmdaction injector.CmdAction
			cmdaction.Cmd = "daftar"
			cmdaction.StepMax = 3
			cmdaction.StepNow = 2
			data.Get("data").ForEach(func(key, value gjson.Result) bool {
				cmdaction.DataPerStep = append(cmdaction.DataPerStep, value.String())
				return true
			})
			cmdaction.DataPerStep = append(cmdaction.DataPerStep, msg.Text)
			di.DelRedisCmd(msg.From.ID)
			json, _ := json.Marshal(cmdaction)
			di.SetRedisCmd(msg.From.ID, string(json), 5*time.Minute)

			rplmsg := di.Config.Get("pesan").Get("daftar").Get("kelompok").String()

			var buttonfin [][]tgbotapi.KeyboardButton
			var button []tgbotapi.KeyboardButton
			col := 0
			di.Config.Get("keyboardlist").Get("kelompok").ForEach(func(key gjson.Result, value gjson.Result) bool {
				var keybuff tgbotapi.KeyboardButton
				keybuff.Text = key.String()
				rplmsg += fmt.Sprintf("\n[%s] %s", key.String(), value.String())
				button = append(button, keybuff)
				col++
				if col == 4 {
					col = 0
					buttonfin = append(buttonfin, button)
					button = nil
				}
				return true
			})
			if cap(button) > 0 {
				buttonfin = append(buttonfin, button)
			}

			var ru tgbotapi.ReplyKeyboardMarkup
			ru.Keyboard = buttonfin
			ru.OneTimeKeyboard = true

			msg := tgbotapi.NewMessage(msg.Chat.ID, rplmsg)
			msg.ReplyMarkup = ru
			di.Enqueue(msg)

		case 2:
			var cmdaction injector.CmdAction
			cmdaction.Cmd = "daftar"
			cmdaction.StepMax = 3
			cmdaction.StepNow = 3
			data.Get("data").ForEach(func(key, value gjson.Result) bool {
				cmdaction.DataPerStep = append(cmdaction.DataPerStep, value.String())
				return true
			})
			cmdaction.DataPerStep = append(cmdaction.DataPerStep, di.Config.Get("keyboardlist").Get("kelompok").Get(msg.Text).String())
			di.DelRedisCmd(msg.From.ID)

			if err := di.TambahRegistrasi(msg.From.ID, cmdaction.DataPerStep[1], cmdaction.DataPerStep[2], cmdaction.DataPerStep[3]); err != nil {
				di.Enqueue(tgbotapi.NewMessage(msg.Chat.ID, di.Config.Get("pesan").Get("daftar").Get("gagal").String()))
				log.Printf("Registrasi: %s\n", err.Error())
			} else {
				di.Enqueue(tgbotapi.NewMessage(msg.Chat.ID, di.Config.Get("pesan").Get("daftar").Get("penutup").String()))
			}

		}

	}
}
