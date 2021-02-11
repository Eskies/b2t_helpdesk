package telebot

import (
	"b2t_helpdesk/injector"
	"log"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/huandu/go-sqlbuilder"
)

func OpenChat(msg *tgbotapi.Message, di *injector.Injector) {
	dbconn, _ := di.DB.Begin()

	var lastticketid int
	if err := dbconn.QueryRow("SELECT MAX(id) FROM tickets WHERE user_id = ?", di.IdTelegramToID(msg.From.ID)).Scan(&lastticketid); err != nil {
		dbconn.Rollback()
		log.Printf("Get Last Ticket: %s\n", err.Error())
		return
	}

	pesan := msg.Text

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

	lokasifile := ""
	//download link
	if len(link) > 0 {
		lokasifile = di.ExPath + "/photos/" + strconv.Itoa(msg.From.ID) + "_"
		if suc := di.DownloadFile(link, &lokasifile); !suc {
			di.Enqueue(tgbotapi.NewMessage(msg.Chat.ID, "[BOT-SYSTEM]\nPesan Anda tidak terkirim, mohon mengirimkan ulang."))
			dbconn.Rollback()
			return
		}
	}

	sqli := sqlbuilder.NewInsertBuilder()
	sqli.InsertInto("chats")
	sqli.Cols("ticket_id", "pesan", "timestamp", "file")
	sqli.Values(lastticketid, pesan, time.Now().Format("2006-01-02 15:04:05"), lokasifile)

	//masukin ke db
	ia, ib := sqli.Build()

	if _, err := dbconn.Exec(ia, ib...); err != nil {
		dbconn.Rollback()
		log.Printf("Add Msg: %s\n", err.Error())
		di.Enqueue(tgbotapi.NewMessage(msg.Chat.ID, "[BOT-SYSTEM]\nPesan Anda tidak terkirim, mohon mengirimkan ulang."))
		return
	}

	dbconn.Commit()

}
