package app

import (
	"github.com/BurntSushi/toml"
	"github.com/rulzurlibrary/api/ext/db"
	"github.com/rulzurlibrary/api/ext/smtp"
	"github.com/rulzurlibrary/api/ext/view"
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
	Paths    struct {
		Static  string
		Thumbs  string
		Cert    string
		Key     string
		Favicon string
	}
}

func ParseConfig() (config Configuration, err error) {
	var base string
	var letsencrypt = path.Join("/", "etc", "letsencrypt", "live", "rulz.xyz")

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
	config.Paths.Cert = path.Join(letsencrypt, "fullchain.pem")
	config.Paths.Key = path.Join(letsencrypt, "privkey.pem")
	config.Paths.Static = path.Join(base, "static")
	config.Paths.Thumbs = path.Join(base, "thumbs")

	config.View.I18n = path.Join(base, "i18n")
	config.View.Templates = path.Join(base, "tplt")
	config.View.Default = "en-US"
	return
}
