package injector

import (
	"context"
	"database/sql"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/adjust/rmq/v3"
	queue "github.com/enriquebris/goconcurrentqueue"
	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/tidwall/gjson"

	_ "github.com/go-sql-driver/mysql"
)

type Injector struct {
	DB        *sql.DB
	Redis     *redis.Client
	MsgOutbox *queue.FIFO
	BotT      *tgbotapi.BotAPI
	WG        sync.WaitGroup
	Ctx       context.Context
	Protokol  string
	Closing   bool
	Config    gjson.Result
	ChannelID int64
	ExPath    string
	OutQ      rmq.Queue
	QConn     rmq.Connection
}

type CmdAction struct {
	Cmd         string   `json:"cmd"`
	StepNow     int      `json:"stepnow"`
	StepMax     int      `json:"stepmax"`
	DataPerStep []string `json:"data"`
}

type QueuePesanKeluar struct {
	ChatId int
	Pesan  string
}

func LoadDependency(lokasikonfigurasi string, expath string) *Injector {
	plan, _ := ioutil.ReadFile(lokasikonfigurasi)
	if !gjson.ValidBytes(plan) {
		log.Panicln("File settings.json tidak terbaca")
	}

	var di Injector
	di.ExPath = expath
	di.Config = gjson.ParseBytes(plan)
	log.Println("Configuration Loaded")

	di.Ctx = context.Background()
	log.Println("Context background")

	redisconf := di.Config.Get("redis")
	di.Redis = redis.NewClient(&redis.Options{
		Addr:     redisconf.Get("host").String() + ":" + redisconf.Get("port").String(),
		Password: redisconf.Get("password").String(),
		DB:       int(redisconf.Get("db").Int()),
	})
	log.Println("Redis Loaded")

	var err error
	dbconf := di.Config.Get("database")
	di.DB, err = sql.Open("mysql", dbconf.Get("user").String()+":"+dbconf.Get("password").String()+"@tcp("+dbconf.Get("host").String()+":"+dbconf.Get("port").String()+")/"+dbconf.Get("db").String())
	if err != nil {
		log.Panicf("Database failed: %s\n", err.Error())
	}
	log.Println("Database Loaded")

	di.BotT, err = tgbotapi.NewBotAPI(di.Config.Get("telegram").Get("botkey").String())
	if err != nil {
		log.Panicf("Database failed: %s\n", err.Error())
	}
	log.Println("Bot Authorized")

	errChan := make(chan error, 10)
	rconn, err := rmq.OpenConnection("outbox-b2t", "tcp", redisconf.Get("host").String()+":"+redisconf.Get("port").String(), int(redisconf.Get("db").Int()), errChan)
	if err != nil {
		log.Panicf("RMQ failed: %s\n", err.Error())
	}
	outQ, err := rconn.OpenQueue("outbox-b2t")
	if err != nil {
		log.Panicf("RMQ failed: %s\n", err.Error())
	}
	if err := outQ.StartConsuming(di.Config.Get("rmq").Get("prefetchlimit").Int(), time.Duration(di.Config.Get("rmq").Get("poolduration_sec").Int())*time.Second); err != nil {
		log.Panicf("RMQ failed: %s\n", err.Error())
	}

	if _, err := outQ.AddBatchConsumer("outbox-b2t", di.Config.Get("rmq").Get("batch_size").Int(), time.Duration(di.Config.Get("rmq").Get("batch_timeout_sec").Int())*time.Second, NewBatchConsumer(&di)); err != nil {
		log.Panicf("RMQ failed: %s\n", err.Error())
	}

	di.OutQ = outQ
	di.QConn = rconn
	go rmqCleaner(&di)
	log.Println("RMQ Loaded")

	di.MsgOutbox = queue.NewFIFO()
	log.Println("Outbox Message Loaded")

	di.ChannelID = di.Config.Get("telegram").Get("channel_id").Int()
	log.Println("Channel ID Loaded")

	return &di
}

func (di *Injector) CloseDependency() {
	di.DB.Close()
	di.Redis.Close()

	<-di.QConn.StopAllConsuming()

	log.Println("Dependency closed.")

	//msg := tgbotapi.NewMessage(di.ChannelID, "Info: Support Bot ditutup. Sampai Jumpa :)")
	//msg.DisableNotification = true
	//di.BotT.Send(msg)
}
