package scrapper

import (
	"github.com/golang/glog"
	"github.com/rulzurlibrary/api/utils"
	"golang.org/x/net/html"
	"net/http"
	"strings"
	"time"
)

const userAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.24 Safari/537.36"

var Client = &http.Client{Timeout: 20 * time.Second}

type ParseFn func(*html.Node) error

func Parse(url string, parseFn ParseFn) error {
	if glog.V(1) {
		glog.Infof(url)
	}
	if strings.Contains(url, "validateCaptcha") {
		return utils.ErrCaptcha
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("User-Agent", userAgent)

	resp, err := Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return err
	}

	return parseFn(doc)
}
