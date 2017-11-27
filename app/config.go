package app

import (
	"github.com/BurntSushi/toml"
	"github.com/RulzUrLibrary/api/ext/db"
	"github.com/RulzUrLibrary/api/ext/i18n"
	"github.com/RulzUrLibrary/api/ext/smtp"
	"github.com/RulzUrLibrary/api/ext/view"
	"os"
	"path"
	"path/filepath"
)

const CONFIG_ENV = "RULZURLIBRARY_CONFIG"

type Configuration struct {
	Debug    bool
	Dev      bool
	Host     string
	Port     int
	Database db.Configuration
	Smtp     smtp.Configuration
	View     view.Configuration
	I18n     i18n.Configuration
	Paths    struct {
		Static   string
		Thumbs   string
		TLSCache string
		Favicon  string
	}
}

func ParseConfig() (config Configuration, err error) {
	var base string

	// get the abs
	// which will try to find the 'filename' from current workind dir too.
	filename, err := filepath.Abs(os.Getenv(CONFIG_ENV))
	if err != nil {
		return
	}

	// put the file's contents as toml to the default configuration(c)
	_, err = toml.DecodeFile(filename, &config)
	if err != nil {
		return
	}

	// infer if we load from current path or from system dir
	if config.Dev {
		base = path.Join(".", "assets")
	} else {
		base = path.Join("/", "var", "lib", "rulzurlibrary")
	}

	// setup various paths
	config.Paths.TLSCache = path.Join(base, "cache")
	config.Paths.Static = path.Join(base, "static")
	config.Paths.Thumbs = path.Join(base, "thumbs")

	config.I18n.Path = path.Join(base, "i18n")
	config.View.Templates = path.Join(base, "tplt")
	return
}
