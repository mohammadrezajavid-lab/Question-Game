package main

import (
	"fmt"
	"golang.project/go-fundamentals/gameapp/scheduler"
	"os"
	"os/signal"
)

func main() {

	done := make(chan bool)
	go func(done <-chan bool) {
		sch := scheduler.New()
		sch.Start(done)
	}(done)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	fmt.Println("received interrupt signal, shutting down gracefully...")

	done <- true
}
