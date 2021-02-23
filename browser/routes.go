package browser

import (
	"b2t_helpdesk/injector"
	"b2t_helpdesk/view"
	"io/ioutil"
	"log"

	"github.com/fasthttp/router"
	"github.com/tidwall/gjson"
	"github.com/valyala/fasthttp"
)

func initRoute(r *router.Router, di *injector.Injector) {
	r.SaveMatchedRoutePath = true
	r.RedirectFixedPath = true
	r.RedirectTrailingSlash = true

	r.ServeFiles("/static/{filepath:*}", di.ExPath+"/assets")
	r.ServeFiles("/photos/{filepath:*}", di.ExPath+"/photos")

	r.GET("/", func(ctx *fasthttp.RequestCtx) {
		p := &view.MainPage{
			CTX: ctx,
		}
		view.WritePageTemplate(ctx, p)
		ctx.SetContentType("text/html; charset=utf-8")
	})

	r.GET("/reload", func(ctx *fasthttp.RequestCtx) {
		plan, _ := ioutil.ReadFile(di.ExPath + "/settings.json")
		if !gjson.ValidBytes(plan) {
			log.Println("Configuration gagal dibaca")
		}

		di.Config = gjson.ParseBytes(plan)
		log.Println("Configuration Loaded")
		ctx.Redirect("/rmq", 200)
	})

	r.GET("/keluhan", func(ctx *fasthttp.RequestCtx) {
		p := &view.KeluhanPage{
			CTX:       ctx,
			Dinjector: di,
		}
		view.WritePageTemplate(ctx, p)
		ctx.SetContentType("text/html; charset=utf-8")
	})

	r.GET("/rmq", func(ctx *fasthttp.RequestCtx) {
		p := &view.RmqPage{
			CTX:       ctx,
			Dinjector: di,
		}
		view.WritePageTemplate(ctx, p)
		ctx.SetContentType("text/html; charset=utf-8")
	})

	r.GET("/jeniskeluhan", func(ctx *fasthttp.RequestCtx) {
		p := &view.JenisKeluhan{
			CTX:       ctx,
			Dinjector: di,
		}
		view.WritePageTemplate(ctx, p)
		ctx.SetContentType("text/html; charset=utf-8")
	})

	//API
	r.POST("/api/datakeluhantable", func(ctx *fasthttp.RequestCtx) {
		dataTableKeluhan(ctx, di)
	})
	r.GET("/api/infoticket/{id}", func(ctx *fasthttp.RequestCtx) {
		infoticket(ctx, di)
	})
	r.GET("/api/chatticket/{id}/{idpesan}", func(ctx *fasthttp.RequestCtx) {
		infochat(ctx, di)
	})

	r.GET("/api/openchat/{id}", func(ctx *fasthttp.RequestCtx) {
		openchat(ctx, di)
	})

	r.GET("/api/closechat/{id}", func(ctx *fasthttp.RequestCtx) {
		closechat(ctx, di)
	})

	r.POST("/api/sendchat/{id}", func(ctx *fasthttp.RequestCtx) {
		sendchat(ctx, di)
	})

	r.GET("/api/openkeluhan/{id}", func(ctx *fasthttp.RequestCtx) {
		openkeluhan(ctx, di)
	})

	r.GET("/api/closekeluhan/{id}", func(ctx *fasthttp.RequestCtx) {
		closekeluhan(ctx, di)
	})
	r.POST("/jeniskeluhan/update", func(ctx *fasthttp.RequestCtx) {
		jenisupdate(ctx, di)
	})
	r.GET("/jeniskeluhan/delete/{id}", func(ctx *fasthttp.RequestCtx) {
		jenisdelete(ctx, di)
	})
	r.GET("/resetdatakeluhan", func(ctx *fasthttp.RequestCtx) {
		di.DB.Exec("TRUNCATE `chats`;")
		di.DB.Exec("TRUNCATE `tickets`;")
		di.DB.Exec("TRUNCATE `users`;")
		ctx.Redirect("/rmq", 200)
	})
	r.GET("/api/infojeniskeluhan/{id}", func(ctx *fasthttp.RequestCtx) {
		type keluhan struct {
			Id        int    `json:"id"`
			Jenis     string `json:"jenis"`
			Hint      string `json:"hint"`
			Autoinput int    `json:"autoinput"`
		}

		var jeniskel keluhan
		id := ctx.UserValue("id").(string)

		_ = di.DB.QueryRow("SELECT id, jenis, hint, autoinput FROM jenisticket WHERE id = ?", id).Scan(&jeniskel.Id, &jeniskel.Jenis, &jeniskel.Hint, &jeniskel.Autoinput)

		di.DoJSONWrite(ctx, 200, jeniskel)
	})
}
