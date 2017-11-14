package main

import (
	"github.com/labstack/gommon/log"
	"github.com/rulzurlibrary/api/app"
	"github.com/rulzurlibrary/api/ext/db"
)

const (
	PREFIX    = "rulz"
	PREFIX_DB = "rulzdb"
)

var (
	level = log.WARN
)

func Logger(prefix string, level log.Lvl) *log.Logger {
	logger := log.New(prefix)
	logger.SetLevel(level)
	return logger
}

func main() {
	config, err := app.ParseConfig()
	logger := log.New(PREFIX)

	if err != nil {
		logger.Fatal(err)
	}

	if config.Debug {
		level = log.DEBUG
	}
	logger.SetLevel(level)

	rulz := app.New(
		db.New(Logger(PREFIX_DB, level), config.Database),
		config,
	)

	// Start application
	rulz.Logger.Fatal(rulz.Start())
}
