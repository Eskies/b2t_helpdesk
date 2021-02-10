package httpserver

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/valyala/fasthttp"
)

func errorCatcher(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		log.Printf("HOST:%s URI:[%s]\n", string(ctx.Host()), string(ctx.RequestURI()))
		next(ctx)
		if ctx.Response.StatusCode() == fasthttp.StatusInternalServerError {
			type errorreport struct {
				Domain   string `json:"domain"`
				Uri      string `json:"uri"`
				Waktu    string `json:"waktu"`
				CTName   string `json:"container"`
				Header   string `json:"header"`
				Request  string `json:"request"`
				Response string `json:"response"`
				ErrorMsg string `json:"errormsg"`
			}

			//error gathering data
			var errordata errorreport
			errordata.Domain = string(ctx.Host())
			errordata.Uri = string(ctx.RequestURI())
			errordata.Waktu = time.Now().Format("02-01-2006 15:04:05")
			errordata.Request = string(ctx.PostBody())
			errordata.Response = string(ctx.Response.Body())
			errordata.ErrorMsg = ctx.UserValue("errormsg").(string)
			errordata.Header = ctx.Request.Header.String()

			log.Printf("HOST:%s URI:[%s]\n\tErr: %s\n", string(ctx.Host()), string(ctx.RequestURI()), errordata.ErrorMsg)

			if ctname, isdock := os.LookupEnv("CTNAME"); isdock {
				errordata.CTName = ctname
				if jdata, err := json.Marshal(errordata); err == nil {
					file, err := os.OpenFile("/errorlog/"+ctname+"-error-log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
					if err == nil {
						defer file.Close()
						file.WriteString(string(jdata))
					} else {
						log.Println(err.Error())
					}
				} else {
					log.Println(err.Error())
				}
			} else {
				errordata.CTName = "develop"
				if jdata, err := json.Marshal(errordata); err == nil {
					file, err := os.OpenFile(errordata.CTName+"-error-log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
					if err == nil {
						defer file.Close()
						file.WriteString(string(jdata))
					} else {
						log.Println(err.Error())
					}
				} else {
					log.Println(err.Error())
				}
			}
		}
	}
}
