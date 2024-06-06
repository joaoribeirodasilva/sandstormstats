package main

import (
	"log/slog"
	"os"

	"github.com/joaoribeirodasilva/sandstormstats/conf"
	"github.com/joaoribeirodasilva/sandstormstats/db"
	"github.com/joaoribeirodasilva/sandstormstats/log_parser"
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

	parser := log_parser.New(db, logger)
	parser.Parse()

}
