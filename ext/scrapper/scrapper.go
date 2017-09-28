package scrapper

import (
	"encoding/base64"
	"github.com/paul-bismuth/library/utils"
	"github.com/labstack/echo"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

const base64Header = "data:image/jpeg;base64,"
const userAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.24 Safari/537.36"

type ParseFn func(*html.Node) error

type Scrapper struct {
	client *http.Client
	logger echo.Logger
	path   string
}

func New(logger echo.Logger, path string) *Scrapper {
	return &Scrapper{
		&http.Client{Timeout: 20 * time.Second},
		logger,
		path,
	}
}

func (s *Scrapper) DownloadThumb(src, isbn string) error {
	return DownloadAsset(src, path.Join(s.path, isbn+".jpg"))
}

func (s *Scrapper) Parse(url string, parseFn ParseFn) error {
	s.logger.Info(url)
	if strings.Contains(url, "validateCaptcha") {
		return utils.ErrCaptcha
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("User-Agent", userAgent)

	resp, err := s.client.Do(req)
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

func DownloadAsset(src, dest string) error {
	if src[:len(base64Header)] == base64Header {
		return base64Decode(src[len(base64Header):], dest)
	} else {
		return downloadAsset(src, dest)
	}
}

func base64Decode(src, dest string) error {
	data, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return err
	}

	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		os.Remove(dest)
		return err
	}
	return nil
}

func downloadAsset(src, dest string) error {
	resp, err := http.Get(src)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		os.Remove(dest)
		return err
	}

	return nil
}
