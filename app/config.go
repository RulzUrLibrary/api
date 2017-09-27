package app

import (
	"github.com/BurntSushi/toml"
	"github.com/ixday/echo-hello/ext/db"
	"path"
	"path/filepath"
)

type Configuration struct {
	Debug    bool             `json:"debug"`
	Dev      bool             `json:"dev"`
	Host     string           `json:"url"`
	Port     int              `json:"port"`
	Database db.Configuration `json:"database"`
	Paths    struct {
		Assets    string
		Thumbs    string
		Templates string
		Cert      string
		Key       string
		Favicon   string
	}
}

func ParseConfig(filename string) (config Configuration, err error) {
	var base string
	var letsencrypt = path.Join("/", "etc", "letsencrypt", "live", "rulz.xyz")

	// get the abs
	// which will try to find the 'filename' from current workind dir too.
	filename, err = filepath.Abs(filename)
	if err != nil {
		return
	}

	// put the file's contents as toml to the default configuration(c)
	_, err = toml.DecodeFile(filename, &config)
	if err != nil {
		return
	}

	if config.Dev {
		base = path.Join(".", "static")
	} else {
		base = path.Join("/", "var", "lib", "rulzurlibrary")
	}
	config.Paths.Cert = path.Join(letsencrypt, "fullchain.pem")
	config.Paths.Key = path.Join(letsencrypt, "privkey.pem")

	return
}
