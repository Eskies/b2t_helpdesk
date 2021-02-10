package injector

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/huandu/go-sqlbuilder"
	"github.com/valyala/fasthttp"
)

func (di *Injector) IsRegistered(userid int) bool {
	var jml int64
	err := di.DB.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", userid).Scan(&jml)
	if err != nil {
		log.Printf("Error isRegistered: %s\n", err.Error())
		return false
	} else {
		if jml == 0 {
			return false
		} else {
			return true
		}
	}
}

func (di *Injector) IsCmdOn(userid int) bool {
	if _, err := di.Redis.Get(di.Ctx, "tbotcmd:"+strconv.Itoa(userid)).Result(); err == nil {
		return true
	} else {
		return false
	}
}

func (di *Injector) CmdData(userid int) string {
	if data, err := di.Redis.Get(di.Ctx, "tbotcmd:"+strconv.Itoa(userid)).Result(); err == nil {
		return data
	} else {
		return ""
	}
}

func (di *Injector) SetRedisCmd(userid int, data string, expire time.Duration) bool {
	if _, err := di.Redis.Set(di.Ctx, "tbotcmd:"+strconv.Itoa(userid), data, expire).Result(); err == nil {
		return true
	} else {
		return false
	}
}

func (di *Injector) DelRedisCmd(userid int) bool {
	if _, err := di.Redis.Del(di.Ctx, "tbotcmd:"+strconv.Itoa(userid)).Result(); err == nil {
		return true
	} else {
		return false
	}
}

func (di *Injector) TambahRegistrasi(userid int, noregis string, nama string, kelompok string) error {
	_, err := di.DB.Exec("INSERT INTO users (id, noregistrasi, nama, kelompok) VALUES (?,?,?,?)", userid, noregis, nama, kelompok)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (di *Injector) AdaBantuan(userid int) bool {

	var cnt int
	err := di.DB.QueryRow("SELECT COUNT(*) FROM tickets WHERE user_id = ? AND close = 0", userid).Scan(&cnt)
	if err != nil || cnt > 0 {
		return true
	} else {
		return false
	}
}

func (di *Injector) IsOpenChat(userid int) bool {

	var cnt int
	err := di.DB.QueryRow("SELECT COUNT(*) FROM users WHERE id = ? AND openchat = 1", userid).Scan(&cnt)
	if err != nil || cnt > 0 {
		return true
	} else {
		return false
	}
}

func (di *Injector) NewTicketBantuan(userid int, jenis string, pesan string, linkfile string) bool {

	dbconn, _ := di.DB.Begin()
	sqli := sqlbuilder.NewInsertBuilder()
	sqli.InsertInto("tickets")
	sqli.Cols("user_id", "jenis")
	sqli.Values(userid, jenis)
	ia, ib := sqli.Build()

	if _, err := dbconn.Exec(ia, ib...); err != nil {
		dbconn.Rollback()
		log.Printf("Add Ticket: %s\n", err.Error())
		return false
	}

	var lastticketid int
	if err := dbconn.QueryRow("SELECT MAX(id) FROM tickets WHERE user_id = ?", userid).Scan(&lastticketid); err != nil {
		dbconn.Rollback()
		log.Printf("Get Last Ticket: %s\n", err.Error())
		return false
	}

	lokasifile := ""

	//download link
	if len(linkfile) > 0 {
		lokasifile = di.ExPath + "/photos/" + strconv.Itoa(userid) + "_"
		if suc := di.DownloadFile(linkfile, &lokasifile); !suc {
			dbconn.Rollback()
			return false
		}
	}

	sqli = nil
	sqli = sqlbuilder.NewInsertBuilder()
	sqli.InsertInto("chats")
	sqli.Cols("ticket_id", "pesan", "timestamp", "file")
	sqli.Values(lastticketid, pesan, time.Now().Format("2006-01-02 15:04:05"), lokasifile)

	//masukin ke db
	ia, ib = sqli.Build()

	if _, err := dbconn.Exec(ia, ib...); err != nil {
		dbconn.Rollback()
		log.Printf("Add Msg: %s\n", err.Error())
		return false
	}

	dbconn.Commit()
	return true
}

func (di *Injector) DownloadFile(fullURL string, location *string) bool {
	// Build fileName from fullPath
	fileURL, err := url.Parse(fullURL)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := *location + segments[len(segments)-1]

	// Create blank file
	file, err := os.Create(fileName)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	// Put content on file
	resp, err := client.Get(fullURL)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	defer resp.Body.Close()

	size, err := io.Copy(file, resp.Body)

	defer file.Close()

	fmt.Printf("Downloaded a file %s with size %d\n", fileName, size)

	segmentsfile := strings.Split(fileName, "/")
	*location = segmentsfile[len(segmentsfile)-1]

	return true
}

func (di *Injector) DoJSONWrite(ctx *fasthttp.RequestCtx, code int, obj interface{}) {
	var (
		strContentType     = []byte("Content-Type")
		strApplicationJSON = []byte("application/json")
	)

	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	ctx.Response.SetStatusCode(code)
	if err := json.NewEncoder(ctx).Encode(obj); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}

func (di *Injector) Enqueue(msg tgbotapi.MessageConfig) {
	taskBytes, err := json.Marshal(msg)
	if err != nil {
		log.Println("Error Marshal")
	}

	err = di.OutQ.PublishBytes(taskBytes)
	if err != nil {
		log.Println("Error Publis Queue")
	}
}
