package telebot

import (
	"b2t_helpdesk/injector"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/tidwall/gjson"
)

func cmdBantuan(msg *tgbotapi.Message, di *injector.Injector, step int) {
	if step == 0 {

		if di.AdaBantuan(msg.From.ID) {
			di.Enqueue(tgbotapi.NewMessage(msg.Chat.ID, di.Config.Get("pesan").Get("bantuan").Get("ongoing").String()))
			return
		}

		rplmsg := di.Config.Get("pesan").Get("bantuan").Get("pembuka").String()

		var buttonfin [][]tgbotapi.KeyboardButton
		var button []tgbotapi.KeyboardButton
		col := 0
		results, _ := di.DB.Query("SELECT id, jenis FROM jenisticket ORDER BY id ASC")
		for results.Next() {
			var id int
			var jenis string

			results.Scan(&id, &jenis)
			var keybuff tgbotapi.KeyboardButton
			keybuff.Text = strconv.Itoa(id)
			rplmsg += fmt.Sprintf("\n[%s] %s", strconv.Itoa(id), jenis)
			button = append(button, keybuff)
			col++
			if col == 5 {
				col = 0
				buttonfin = append(buttonfin, button)
				button = nil
			}

		}

		di.Config.Get("keyboardlist").Get("jenisticket").ForEach(func(key gjson.Result, value gjson.Result) bool {
			var keybuff tgbotapi.KeyboardButton
			keybuff.Text = key.String()
			rplmsg += fmt.Sprintf("\n[%s] %s", key.String(), value.Get("jenis").String())
			button = append(button, keybuff)
			col++
			if col == 5 {
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

		msgo := tgbotapi.NewMessage(msg.Chat.ID, rplmsg)
		msgo.ReplyMarkup = ru
		di.Enqueue(msgo)

		var cmdaction injector.CmdAction
		cmdaction.Cmd = "bantuan"
		cmdaction.StepMax = 99
		cmdaction.StepNow = 0
		cmdaction.DataPerStep = append(cmdaction.DataPerStep, "start")

		json, _ := json.Marshal(cmdaction)

		di.SetRedisCmd(msg.From.ID, string(json), 5*time.Minute)
	} else {
		data := gjson.Parse(di.CmdData(msg.From.ID))
		switch gjson.Parse(di.CmdData(msg.From.ID)).Get("stepnow").Int() {
		case 0:
			var autoinput bool
			var hint string
			var jenisticket string

			if strings.Contains(msg.Text, "?") {
				hint = di.Config.Get("keyboardlist").Get("jenisticket").Get(msg.Text).Get("hint").String()
				jenisticket = di.Config.Get("keyboardlist").Get("jenisticket").Get(msg.Text).Get("jenis").String()
				autoinput = di.Config.Get("keyboardlist").Get("jenisticket").Get(msg.Text).Get("autoinput").Bool()
			} else {
				var autoinput_buf int
				di.DB.QueryRow("SELECT hint, jenis, autoinput FROM jenisticket WHERE id = ?", msg.Text).Scan(&hint, &jenisticket, &autoinput_buf)
				if autoinput_buf == 0 {
					autoinput = false
				} else {
					autoinput = true
				}
			}

			if !autoinput {
				var ru tgbotapi.ReplyKeyboardMarkup
				ru.Keyboard = keyboardYaTidak()
				ru.OneTimeKeyboard = true

				msgo := tgbotapi.NewMessage(msg.Chat.ID, hint+"\n\n"+di.Config.Get("pesan").Get("bantuan").Get("konfirmasi").String())
				msgo.ReplyMarkup = ru
				di.Enqueue(msgo)

				var cmdaction injector.CmdAction
				cmdaction.Cmd = "bantuan"
				cmdaction.StepMax = 99
				cmdaction.StepNow = 1
				cmdaction.DataPerStep = append(cmdaction.DataPerStep, jenisticket)

				json, _ := json.Marshal(cmdaction)

				di.DelRedisCmd(msg.From.ID)
				di.SetRedisCmd(msg.From.ID, string(json), 5*time.Minute)
			} else {
				di.Enqueue(tgbotapi.NewMessage(msg.Chat.ID, hint+"\n\n"+di.Config.Get("pesan").Get("bantuan").Get("sebelummodeinput").String()))

				var cmdaction injector.CmdAction
				cmdaction.Cmd = "bantuan"
				cmdaction.StepMax = 99
				cmdaction.StepNow = 2
				cmdaction.DataPerStep = append(cmdaction.DataPerStep, jenisticket)

				json, _ := json.Marshal(cmdaction)

				di.DelRedisCmd(msg.From.ID)
				di.SetRedisCmd(msg.From.ID, string(json), 10*time.Minute)
			}
		case 1:
			if strings.ToLower(msg.Text) == "iya" {
				di.Enqueue(tgbotapi.NewMessage(msg.Chat.ID, di.Config.Get("pesan").Get("bantuan").Get("sebelummodeinput").String()))

				var cmdaction injector.CmdAction
				cmdaction.Cmd = "bantuan"
				cmdaction.StepMax = 99
				cmdaction.StepNow = 2
				data.Get("data").ForEach(func(key, value gjson.Result) bool {
					cmdaction.DataPerStep = append(cmdaction.DataPerStep, value.String())
					return true
				})

				json, _ := json.Marshal(cmdaction)

				di.DelRedisCmd(msg.From.ID)
				di.SetRedisCmd(msg.From.ID, string(json), 10*time.Minute)
			} else {
				di.DelRedisCmd(msg.From.ID)
				di.Enqueue(tgbotapi.NewMessage(msg.Chat.ID, di.Config.Get("pesan").Get("bantuan").Get("selesai").String()))
			}

		case 2:
			var cmdaction injector.CmdAction
			cmdaction.Cmd = "bantuan"
			cmdaction.StepMax = 99
			cmdaction.StepNow = 3
			data.Get("data").ForEach(func(key, value gjson.Result) bool {
				cmdaction.DataPerStep = append(cmdaction.DataPerStep, value.String())
				return true
			})
			cmdaction.DataPerStep = append(cmdaction.DataPerStep, msg.Text)

			json, _ := json.Marshal(cmdaction)

			di.DelRedisCmd(msg.From.ID)
			di.SetRedisCmd(msg.From.ID, string(json), 10*time.Minute)

			di.Enqueue(tgbotapi.NewMessage(msg.Chat.ID, di.Config.Get("pesan").Get("bantuan").Get("sebelummodegambar").String()))
		case 3:
			var cmdaction injector.CmdAction
			cmdaction.Cmd = "bantuan"
			cmdaction.StepMax = 3
			cmdaction.StepNow = 1
			data.Get("data").ForEach(func(key, value gjson.Result) bool {
				cmdaction.DataPerStep = append(cmdaction.DataPerStep, value.String())
				return true
			})

			maxsize := 0
			link := ""
			if msg.Photo != nil {
				for _, pht := range *msg.Photo {
					linkfile, err := di.BotT.GetFileDirectURL(pht.FileID)
					if err == nil {
						if maxsize < pht.FileSize {
							maxsize = pht.FileSize
							link = linkfile
						}
					}
				}
			}
			var pesanakhir string
			if di.NewTicketBantuan(msg.From.ID, cmdaction.DataPerStep[0], cmdaction.DataPerStep[1], link) {
				pesanakhir = di.Config.Get("pesan").Get("bantuan").Get("penutup").String()
			} else {
				pesanakhir = di.Config.Get("pesan").Get("bantuan").Get("gagal").String()
			}
			di.Enqueue(tgbotapi.NewMessage(msg.Chat.ID, pesanakhir+"\n\n"+di.Config.Get("pesan").Get("bantuan").Get("selesai").String()))
			di.DelRedisCmd(msg.From.ID)
		}

	}
}
