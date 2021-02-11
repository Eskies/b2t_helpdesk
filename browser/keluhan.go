package browser

import (
	"b2t_helpdesk/injector"
	"database/sql"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jinzhu/copier"
	"github.com/valyala/fasthttp"
)

func dataTableKeluhan(ctx *fasthttp.RequestCtx, di *injector.Injector) {
	dbConn := di.DB

	draw, _ := strconv.Atoi(string(ctx.FormValue("draw")))
	start, _ := strconv.Atoi(string(ctx.FormValue("start")))
	length, _ := strconv.Atoi(string(ctx.FormValue("length")))
	search := string(ctx.FormValue("search[value]"))
	ordercol, _ := strconv.Atoi(string(ctx.FormValue("order[0][column]")))
	orderdir := string(ctx.FormValue("order[0][dir]"))

	filterjenis := string(ctx.FormValue("jenisfilter"))
	filterkelompok := string(ctx.FormValue("kelompokfilter"))
	filterstatus := string(ctx.FormValue("statusfilter"))
	filteropen := string(ctx.FormValue("openfilter"))

	//List Search Coloumn
	orderColomns := []string{
		"tickets.id",
		"users.noregistrasi",
		"users.nama",
		"users.kelompok",
		"tickets.jenis",
		"tickets.close",
		"users.openchat",
		"tickets.tim",
	}

	//tipe struct reply
	type datastruct struct {
		ID       int    `json:"id"`
		NoReg    string `json:"no"`
		Nama     string `json:"nama"`
		Kelompok string `json:"kelompok"`
		Jenis    string `json:"jenis"`
		Close    string `json:"status"`
		Openchat string `json:"openchat"`
		Tim      string `json:"tim"`
	}
	type datatablestruct struct {
		Draw            int          `json:"draw"`
		RecordsTotal    int          `json:"recordsTotal"`
		RecordsFiltered int          `json:"recordsFiltered"`
		Data            []datastruct `json:"data"`
		Debug           string       `json:"debug"`
	}

	var dtreply datatablestruct
	dtreply.Draw = draw
	dtreply.Data = []datastruct{}

	sql := sqlbuilder.NewSelectBuilder()
	sql.From("tickets")
	sql.Join("users", "users.id = tickets.user_id")

	//count all
	sqlC := sql
	copier.Copy(sqlC, sql)
	sqlC.Select("COUNT(*)")

	//get count all
	cA, cB := sqlC.Build()
	err := dbConn.QueryRow(cA, cB...).Scan(&dtreply.RecordsTotal)
	if err != nil {
		_, _ = ctx.WriteString("Internal Error, Contact Admin!")
		ctx.SetUserValue("errormsg", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetConnectionClose()
		return
	}

	//filter
	if filterjenis != "-1" {
		sql.Where(sql.E("tickets.jenis", filterjenis))
	}
	if filterkelompok != "-1" {
		sql.Where(sql.E("users.kelompok", filterkelompok))
	}
	if filterstatus != "-1" {
		sql.Where(sql.E("tickets.close", filterstatus))
	}
	if filteropen != "-1" {
		sql.Where(sql.E("users.openchat", filteropen))
	}

	//Search Keyword
	if len(search) > 0 {
		keyword := sqlbuilder.Named("keyword", "%"+search+"%")
		sql.Where(
			sql.Or(
				sql.Like("users.noregistrasi", keyword),
				sql.Like("users.nama", keyword),
				sql.Like("users.kelompok", keyword),
				sql.Like("tickets.jenis", keyword),
				sql.Like("tickets.tim", keyword),
			),
		)
	}

	//count filtered
	sqlF := sql
	copier.Copy(sqlF, sql)
	sqlF.Select("COUNT(*)")

	//get filtered count
	aF, bF := sqlF.Build()
	err = dbConn.QueryRow(aF, bF...).Scan(&dtreply.RecordsFiltered)
	if err != nil {
		_, _ = ctx.WriteString("Internal Error, Contact Admin!")
		ctx.SetUserValue("errormsg", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetConnectionClose()
		return
	}

	//Order dir
	sql.Select(
		"tickets.id",
		"users.noregistrasi",
		"users.nama",
		"users.kelompok",
		"tickets.jenis",
		"IF("+sql.E("tickets.close", 0)+", 'TERBUKA', 'TERTUTUP') as close",
		"IF("+sql.E("users.openchat", 0)+", 'CLOSE', 'OPEN') as openchat",
		"COALESCE(tickets.tim, '-')",
	)

	sql.OrderBy(orderColomns[ordercol] + " " + orderdir)
	sql.Limit(length)
	sql.Offset(start)

	//get results
	a, b := sql.Build()
	qresults, err := dbConn.Query(a, b...)
	if err != nil {
		_, _ = ctx.WriteString("Internal Error, Contact Admin!")
		ctx.SetUserValue("errormsg", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetConnectionClose()
		return
	}

	//dtreply.DebugSQL = a
	defer qresults.Close()
	for qresults.Next() {
		var buffData datastruct
		err = qresults.Scan(
			&buffData.ID,
			&buffData.NoReg,
			&buffData.Nama,
			&buffData.Kelompok,
			&buffData.Jenis,
			&buffData.Close,
			&buffData.Openchat,
			&buffData.Tim,
		)

		if err != nil {
			_, _ = ctx.WriteString("Internal Error, Contact Admin!")
			ctx.SetUserValue("errormsg", err.Error())
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.SetConnectionClose()
			return
		}

		dtreply.Data = append(dtreply.Data, buffData)
	}
	di.DoJSONWrite(ctx, 200, dtreply)
}

func infoticket(ctx *fasthttp.RequestCtx, di *injector.Injector) {
	idticket := ctx.UserValue("id").(string)
	type datastruct struct {
		ID       int    `json:"id"`
		NoReg    string `json:"no"`
		Nama     string `json:"nama"`
		Kelompok string `json:"kelompok"`
		Jenis    string `json:"jenis"`
		Close    string `json:"status"`
		Openchat string `json:"openchat"`
	}

	var reply datastruct

	err := di.DB.QueryRow(
		`SELECT 
			users.nama,
			users.noregistrasi, 
			users.kelompok, 
			tickets.jenis,
			tickets.close,
			users.openchat
		FROM tickets
		INNER JOIN users ON users.id = tickets.user_id
		WHERE tickets.id = ?`,
		idticket,
	).Scan(
		&reply.Nama,
		&reply.NoReg,
		&reply.Kelompok,
		&reply.Jenis,
		&reply.Close,
		&reply.Openchat,
	)

	if err != nil {
		_, _ = ctx.WriteString("Internal Error, Contact Admin!")
		ctx.SetUserValue("errormsg", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetConnectionClose()
		return
	}

	di.DoJSONWrite(ctx, 200, reply)
}

func infochat(ctx *fasthttp.RequestCtx, di *injector.Injector) {
	idticket := ctx.UserValue("id").(string)
	idpesan := ctx.UserValue("idpesan").(string)

	type datastruct struct {
		ID        int    `json:"id"`
		Pesan     string `json:"pesan"`
		File      string `json:"file"`
		Timestamp string `json:"timestamp"`
		Penulis   string `json:"penulis"`
	}

	type dataset struct {
		Chat []datastruct `json:"chat"`
	}

	var reply dataset

	results, err := di.DB.Query(
		`SELECT 
			chats.id,
			chats.pesan,
			COALESCE(chats.file, ''),
			DATE_FORMAT(chats.timestamp, '%d-%m-%Y %H:%i:%s'),
			COALESCE(chats.penulis, '')
		FROM chats
		WHERE chats.ticket_id = ? AND chats.id > ?
		ORDER BY chats.id ASC`,
		idticket, idpesan,
	)

	if err != nil {
		_, _ = ctx.WriteString("Internal Error, Contact Admin!")
		ctx.SetUserValue("errormsg", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetConnectionClose()
		return
	} else {
		for results.Next() {
			var replyx datastruct

			err := results.Scan(
				&replyx.ID,
				&replyx.Pesan,
				&replyx.File,
				&replyx.Timestamp,
				&replyx.Penulis,
			)
			if err != nil {
				_, _ = ctx.WriteString("Internal Error, Contact Admin!")
				ctx.SetUserValue("errormsg", err.Error())
				ctx.SetStatusCode(fasthttp.StatusInternalServerError)
				ctx.SetConnectionClose()
				return
			} else {
				replyx.Pesan = strings.Replace(replyx.Pesan, "\n", "<br>", -1)
				reply.Chat = append(reply.Chat, replyx)
			}
		}
	}

	di.DoJSONWrite(ctx, 200, reply)
}

func openchat(ctx *fasthttp.RequestCtx, di *injector.Injector) {
	idticket := ctx.UserValue("id").(string)

	_, err := di.DB.Exec("UPDATE users INNER JOIN tickets ON tickets.user_id = users.id SET users.openchat = 1 WHERE tickets.id = ?", idticket)
	if err != nil {
		_, _ = ctx.WriteString("Internal Error, Contact Admin!")
		ctx.SetUserValue("errormsg", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetConnectionClose()
		return
	}

	var userid int64
	err = di.DB.QueryRow("SELECT users.idtelegram FROM tickets INNER JOIN users ON users.id = tickets.user_id WHERE tickets.id = ?", idticket).Scan(&userid)
	if err != nil {
		_, _ = ctx.WriteString("Internal Error, Contact Admin!")
		ctx.SetUserValue("errormsg", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetConnectionClose()
		return
	}

	di.Enqueue(tgbotapi.NewMessage(userid, di.Config.Get("pesan").Get("openchat").Get("open").String()))

}

func closechat(ctx *fasthttp.RequestCtx, di *injector.Injector) {
	idticket := ctx.UserValue("id").(string)

	_, err := di.DB.Exec("UPDATE users INNER JOIN tickets ON tickets.user_id = users.id SET users.openchat = 0 WHERE tickets.id = ?", idticket)
	if err != nil {
		_, _ = ctx.WriteString("Internal Error, Contact Admin!")
		ctx.SetUserValue("errormsg", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetConnectionClose()
		return
	}

	var userid int64
	err = di.DB.QueryRow("SELECT users.idtelegram FROM tickets INNER JOIN users ON users.id = tickets.user_id WHERE tickets.id = ?", idticket).Scan(&userid)
	if err != nil {
		_, _ = ctx.WriteString("Internal Error, Contact Admin!")
		ctx.SetUserValue("errormsg", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetConnectionClose()
		return
	}

	di.Enqueue(tgbotapi.NewMessage(userid, di.Config.Get("pesan").Get("openchat").Get("close").String()))

}

func sendchat(ctx *fasthttp.RequestCtx, di *injector.Injector) {
	idticket := ctx.UserValue("id").(string)
	pesan := string(ctx.FormValue("pesan"))
	pegawai := string(ctx.FormValue("pegawai"))

	var userid int64
	err := di.DB.QueryRow("SELECT users.idtelegram FROM users INNER JOIN tickets ON tickets.user_id = users.id WHERE tickets.id = ? AND users.deleted_at IS NULL;", idticket).Scan(&userid)
	if err != nil {
		if err == sql.ErrNoRows {
			_, _ = ctx.WriteString("Pengguna sudah tidak terdaftar!")
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		} else {
			_, _ = ctx.WriteString("Internal Error, Contact Admin!")
			ctx.SetUserValue("errormsg", err.Error())
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.SetConnectionClose()
			return
		}
	}

	sqli := sqlbuilder.NewInsertBuilder()
	sqli.InsertInto("chats")
	sqli.Cols("ticket_id", "pesan", "timestamp", "penulis")
	sqli.Values(idticket, pesan, time.Now().Format("2006-01-02 15:04:05"), "TIM-"+strings.ToLower(pegawai))

	//masukin ke db
	ia, ib := sqli.Build()

	if _, err := di.DB.Exec(ia, ib...); err != nil {
		_, _ = ctx.WriteString("Internal Error, Contact Admin!")
		ctx.SetUserValue("errormsg", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetConnectionClose()
		return
	}

	di.DB.Exec("UPDATE tickets SET tim = ? WHERE id = ?", strings.ToLower(pegawai), idticket)

	di.Enqueue(tgbotapi.NewMessage(userid, "[Dari TIM-Support]\n"+pesan))
}

func openkeluhan(ctx *fasthttp.RequestCtx, di *injector.Injector) {
	idticket := ctx.UserValue("id").(string)

	_, err := di.DB.Exec("UPDATE tickets SET close = 0 WHERE tickets.id = ?", idticket)
	if err != nil {
		_, _ = ctx.WriteString("Internal Error, Contact Admin!")
		ctx.SetUserValue("errormsg", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetConnectionClose()
		return
	}

	var userid int64
	err = di.DB.QueryRow("SELECT users.idtelegram FROM users INNER JOIN tickets ON tickets.user_id = users.id WHERE tickets.id = ? AND users.deleted_at IS NULL;", idticket).Scan(&userid)
	if err != nil {
		_, _ = ctx.WriteString("Internal Error, Contact Admin!")
		ctx.SetUserValue("errormsg", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetConnectionClose()
		return
	}

	di.Enqueue(tgbotapi.NewMessage(userid, di.Config.Get("pesan").Get("bantuan").Get("open").String()))
}

func closekeluhan(ctx *fasthttp.RequestCtx, di *injector.Injector) {
	idticket := ctx.UserValue("id").(string)

	_, err := di.DB.Exec("UPDATE tickets SET close = 1 WHERE tickets.id = ?", idticket)
	if err != nil {
		_, _ = ctx.WriteString("Internal Error, Contact Admin!")
		ctx.SetUserValue("errormsg", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetConnectionClose()
		return
	}

	var userid int64
	err = di.DB.QueryRow("SELECT users.idtelegram FROM users INNER JOIN tickets ON tickets.user_id = users.id WHERE tickets.id = ? AND users.deleted_at IS NULL;", idticket).Scan(&userid)
	if err != nil {
		_, _ = ctx.WriteString("Internal Error, Contact Admin!")
		ctx.SetUserValue("errormsg", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetConnectionClose()
		return
	}

	di.Enqueue(tgbotapi.NewMessage(userid, di.Config.Get("pesan").Get("bantuan").Get("close").String()))

}
