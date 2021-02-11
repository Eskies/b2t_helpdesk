package main

import (
	"b2t_helpdesk/browser"
	"b2t_helpdesk/injector"
	"b2t_helpdesk/telebot"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func main() {
	ex, err := os.Executable()
	if err != nil {
		log.Panicln(err.Error())
	}
	exPath := filepath.Dir(ex)
	dinjector := injector.LoadDependency(exPath+"/settings.json", exPath)

	dinjector.WG.Add(1)
	dinjector.Closing = false

	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Dependency Loaded!")

	go telebot.TelebotStart(dinjector)
	go browser.StartHTTPServer(dinjector)

	<-termChan

	dinjector.Closing = true
	dinjector.CloseDependency()
	dinjector.WG.Wait()

	log.Println("Bye!")
}
