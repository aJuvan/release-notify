package main

import (
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aJuvan/release-notify/config"
	"github.com/aJuvan/release-notify/notify"
)

func main() {
	godotenv.Load()

	conf, err := config.ParseConfig()
	if err != nil {
		fmt.Println("Error parsing config file")
		panic(err)
	}

	notify.Setup(conf)

	loc, err := time.LoadLocation(conf.Core.Location)
	if err != nil {
		fmt.Println("Error loading location")
		panic(err)
	}

	s, err := gocron.NewScheduler(gocron.WithLocation(loc))
	if err != nil {
		fmt.Println("Error creating scheduler")
		panic(err)
	}

	s.NewJob(
		gocron.CronJob(conf.Core.Cron, false),
		gocron.NewTask(func() {
			notify.Trigger()
		}),
	)

	s.Start()
	defer s.Shutdown()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	select {
	case <-c:
	}
}
