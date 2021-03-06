package injector

import (
	"encoding/json"
	"log"
	"time"

	"github.com/adjust/rmq/v3"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type BatchConsumer struct {
	di *Injector
}

func NewBatchConsumer(di *Injector) *BatchConsumer {
	return &BatchConsumer{di: di}
}

func (consumer *BatchConsumer) Consume(batch rmq.Deliveries) {
	payloads := batch.Payloads()

	var msg tgbotapi.MessageConfig
	for _, payload := range payloads {
		if err := json.Unmarshal([]byte(payload), &msg); err != nil {
			log.Println("Fail unmarshal msg")
		}
		consumer.di.BotT.Send(msg)
	}
	errors := batch.Ack()
	if len(errors) == 0 {
		return
	}
	time.Sleep(1 * time.Second)
}

func rmqCleaner(di *Injector) {
	cleaner := rmq.NewCleaner(di.QConn)

	for range time.Tick(time.Minute) {
		_, err := cleaner.Clean()
		if err != nil {
			log.Printf("RMQ Cleaner: failed to clean: %s\n", err)
			continue
		}

		_, err = di.OutQ.ReturnRejected(100)
		if err != nil {
			log.Printf("RMQ Cleaner: return rejected failed: %s\n", err)
			continue
		}
	}
}
