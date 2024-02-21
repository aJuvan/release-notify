package main

import (
	"github.com/go-co-op/gocron/v2"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/aJuvan/release-notify/config"
	"github.com/aJuvan/release-notify/notify"
)

func main() {
	// Load configuration
	godotenv.Load()

	conf, err := config.ParseConfig()
	if err != nil {
		log.Error().Err(err).Msg("Error parsing config file")
		panic(err)
	}

	log.Debug().Msgf("Loaded Config %+v", conf)

	// Setup notification channels and releases
	notify.Setup(conf)

	// Setup scheduler
	loc, err := time.LoadLocation(conf.Core.Location)
	if err != nil {
		log.Error().Err(err).Msg("Error loading location")
		panic(err)
	}

	log.Debug().Msgf("Using location: %s", loc)

	s, err := gocron.NewScheduler(gocron.WithLocation(loc))
	if err != nil {
		log.Error().Err(err).Msg("Error creating scheduler")
		panic(err)
	}

	// Schedule job
	s.NewJob(
		gocron.CronJob(conf.Core.Cron, false),
		gocron.NewTask(func() {
			log.Debug().Msg("Running trigger")
			notify.Trigger()
			log.Debug().Msg("Trigger finished")
		}),
	)

	// Start scheduler and wait for interrupt
	s.Start()
	defer s.Shutdown()

	log.Info().Msg("Scheduler started")

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	select {
	case <-c:
	}
}
