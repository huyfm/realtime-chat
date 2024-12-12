package main

import (
	"context"
	"os"
	"time"

	"github.com/huyfm/rtc"
	"github.com/huyfm/rtc/http"
	"github.com/huyfm/rtc/pg"
)

func main() {
	conf, err := rtc.DefaultConfig()
	if err != nil {
		rtc.Logger.Error().Msg("can't read .env")
		os.Exit(1)
	}

	ctx := context.Background()
	db, err := pg.OpenDB(ctx, conf.DSN)
	if err != nil {
		rtc.Logger.Fatal().Err(err).Msg("can't connect to database")
		os.Exit(1)
	}
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	if err := db.Ping(ctx); err != nil {
		rtc.Logger.Fatal().Err(err).Msg("can't connect to database")
		os.Exit(1)
	}
	defer pg.CloseDB(db)

	s := http.NewServer(conf)
	// Assign services manually.
	s.UserSrv = pg.NewUserService(db)

	if err := s.Open(); err != nil {
		rtc.Logger.Fatal().Err(err).Msg("server error")
		os.Exit(1)
	}

	// TODO: gracefully shutdown.
}
