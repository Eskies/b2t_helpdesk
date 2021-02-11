package browser

import (
	"b2t_helpdesk/injector"

	"github.com/valyala/fasthttp"
)

func jenisupdate(ctx *fasthttp.RequestCtx, di *injector.Injector) {
	id := string(ctx.FormValue("id"))
	jenis := string(ctx.FormValue("jenis"))
	hint := string(ctx.FormValue("hint"))
	autoinput := string(ctx.FormValue("autoinput"))

	if id == "0" {
		di.DB.Exec("INSERT INTO jenisticket (jenis, hint, autoinput) VALUES (?,?,?)", jenis, hint, autoinput)
	} else {
		di.DB.Exec("UPDATE jenisticket SET jenis = ?, hint = ?, autoinput = ? WHERE id = ?", jenis, hint, autoinput, id)
	}

	ctx.Redirect("/jeniskeluhan", 200)
}

func jenisdelete(ctx *fasthttp.RequestCtx, di *injector.Injector) {
	idkeluhan := ctx.UserValue("id").(string)
	di.DB.Exec("DELETE FROM jenisticket WHERE id = ?", idkeluhan)
	ctx.Redirect("/jeniskeluhan", 200)
}
