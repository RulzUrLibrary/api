package app

import (
	"github.com/labstack/gommon/log"
	"github.com/rulzurlibrary/api/ext/auth"
	"github.com/rulzurlibrary/api/ext/db"
	"github.com/rulzurlibrary/api/ext/scrapper"
	"github.com/rulzurlibrary/api/ext/smtp"
	"github.com/rulzurlibrary/api/ext/validator"
	"github.com/rulzurlibrary/api/ext/view"
)

const (
	PREFIX          = "rulz"
	PREFIX_DB       = "rulz_db"
	PREFIX_CACHE    = "rulz_cache"
	PREFIX_AUTH     = "rulz_auth"
	PREFIX_SMTP     = "rulz_smtp"
	PREFIX_SCRAPPER = "rulz_scrapper"
)

type Initializer interface {
	DB() (*db.DB, *auth.Auth)
	Smtp() *smtp.Smtp
	View(*Application) *view.View
	Scrapper() *scrapper.Scrapper
	Validator() *validator.Validator
	Config() Configuration
	Logger(prefix string) *log.Logger
}

type DefaultInitializer struct {
	Configuration
	log.Lvl
}

func NewDefaultInitializer() *DefaultInitializer {
	lvl := log.WARN // default to WARN
	config, err := ParseConfig()

	if err != nil {
		log.New(PREFIX).Fatal(err)
	}

	if config.Debug {
		lvl = log.DEBUG
	}
	return &DefaultInitializer{config, lvl}
}

func (di *DefaultInitializer) Logger(prefix string) *log.Logger {
	logger := log.New(prefix)
	logger.SetLevel(di.Lvl)
	return logger
}

func (di *DefaultInitializer) Smtp() *smtp.Smtp {
	return smtp.New(di.Logger(PREFIX_SMTP), di.Configuration.Smtp)
}

func (di *DefaultInitializer) View(app *Application) *view.View {
	return view.New(app.Web, di.Configuration.View)
}

func (di *DefaultInitializer) Scrapper() *scrapper.Scrapper {
	return scrapper.New(di.Logger(PREFIX_SCRAPPER), di.Configuration.Paths.Thumbs)
}

func (di *DefaultInitializer) Validator() *validator.Validator {
	return validator.New()
}

func (di *DefaultInitializer) Config() Configuration {
	return di.Configuration
}

func (di *DefaultInitializer) DB() (*db.DB, *auth.Auth) {
	database := db.New(di.Logger(PREFIX_DB), di.Configuration.Database)
	auth := auth.New(di.Logger(PREFIX_AUTH), database, auth.NewDefaultCache(di.Logger(PREFIX_CACHE)))
	return database, auth
}
