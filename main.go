package main

import (
	"log/slog"
	"os"

	"github.com/joaoribeirodasilva/sandstormstats/conf"
	"github.com/joaoribeirodasilva/sandstormstats/db"
)

func main() {

	logOpts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, logOpts))
	slog.SetDefault(logger)

	conf := conf.New(logger)
	conf.Read()

	db := db.New(conf.GetDbConf(), logger)
	db.Connect()

	//LogRead("./log/2f70c02c-a71d-4404-bb70-d486d0692f49-backup-2024.06.04-16.02.16.log")
}
