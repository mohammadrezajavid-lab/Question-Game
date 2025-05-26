package main

import (
	"fmt"
	"golang.project/go-fundamentals/gameapp/config/httpservercfg"
	"golang.project/go-fundamentals/gameapp/config/setupservices"
	"golang.project/go-fundamentals/gameapp/scheduler"
	"os"
	"os/signal"
	"sync"
)

func main() {

	// setup config
	config := httpservercfg.NewConfig("", 0)
	fmt.Println("config project: ", config)

	// setup services
	setupSvc := setupservices.New(config)

	// start scheduler goroutine
	var wg sync.WaitGroup
	done := make(chan bool)
	sch := scheduler.New(setupSvc.MatchingSvc)
	wg.Add(1)
	go sch.Start(done, &wg)

	// waiting for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit // blocked main in this line
	fmt.Println("received interrupt signal, shutting down gracefully...")

	// stopping scheduler
	done <- true
	wg.Wait()
}
