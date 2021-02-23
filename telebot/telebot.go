package telebot

import (
	"b2t_helpdesk/injector"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/tidwall/gjson"
)

func TelebotStart(di *injector.Injector) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := di.BotT.GetUpdatesChan(u)

	if err != nil {
		log.Fatalf("Chan Bot Failed: %s\n", err.Error())
	}

	//msg := tgbotapi.NewMessage(di.ChannelID, "Info: Support Bot telah dibuka. Silahkan chat kami untuk mengajukan bantuan.")
	//msg.DisableNotification = true
	//di.BotT.Send(msg)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		go prosesPesanMasuk(update.Message, di)

	}
}

func prosesPesanMasuk(msg *tgbotapi.Message, di *injector.Injector) {
	if di.Config.Get("telegram").Get("close").Bool() {
		msg := tgbotapi.NewMessage(msg.Chat.ID, di.Config.Get("telegram").Get("pesanclose").String())
		di.Enqueue(msg)
		return
	}

	//cek jam operasional
	jambuka, _ := time.Parse("2006-01-02 15:04:05 Z0700", time.Now().Format("2006-01-02 ")+di.Config.Get("telegram").Get("jamaktif").Get("mulai").String())
	jamtutup, _ := time.Parse("2006-01-02 15:04:05 Z0700", time.Now().Format("2006-01-02 ")+di.Config.Get("telegram").Get("jamaktif").Get("tutup").String())
	loc, _ := time.LoadLocation("Asia/Makassar")
	harikey := time.Now().In(loc).Format("Mon")
	if time.Now().Before(jambuka) || time.Now().After(jamtutup) || di.Config.Get("telegram").Get("jamaktif").Get("hariaktif").Get(harikey).Bool() == false {
		msg := tgbotapi.NewMessage(msg.Chat.ID, di.Config.Get("telegram").Get("jamaktif").Get("pesantutup").String())
		di.Enqueue(msg)
		return
	}

	if di.IsOpenChat(msg.From.ID) {
		OpenChat(msg, di)
		return
	}
	if strings.ToLower(msg.Command()) == "batal" {
		di.DelRedisCmd(msg.From.ID)
		msg := tgbotapi.NewMessage(msg.Chat.ID, di.Config.Get("pesan").Get("batal").String())
		di.Enqueue(msg)
	} else {
		if di.IsRegistered(msg.From.ID) {
			if !di.IsCmdOn(msg.From.ID) {
				switch strings.ToLower(msg.Command()) {
				case "cmd":
					msg := tgbotapi.NewMessage(msg.Chat.ID, di.Config.Get("pesan").Get("cmd").String())
					di.Enqueue(msg)
				case "bantuan":
					cmdBantuan(msg, di, 0)
				case "unreg":
					deleteRegistrasi(msg, di)
				default:
					msg := tgbotapi.NewMessage(msg.Chat.ID, di.Config.Get("pesan").Get("default").String())
					di.Enqueue(msg)
				}
			} else {
				switch gjson.Parse(di.CmdData(msg.From.ID)).Get("cmd").String() {
				case "bantuan":
					cmdBantuan(msg, di, 1)
				}
			}
		} else {
			if strings.ToLower(msg.Command()) == "daftar" {
				cmdDaftar(msg, di, true)
			} else {
				if di.IsCmdOn(msg.From.ID) {
					cmdDaftar(msg, di, false)
				} else {
					msg := tgbotapi.NewMessage(msg.Chat.ID, di.Config.Get("pesan").Get("welcome").String())
					di.Enqueue(msg)
				}

			}

		}
	}
}
