package scrapper

import (
	"encoding/base64"
	"github.com/rulzurlibrary/api/ext/db"
	"github.com/rulzurlibrary/api/utils"
	"io"
	"net/http"
	"os"
	"path"
)

const base64Header = "data:image/jpeg;base64,"

type Book struct {
	*utils.Book
}

func (b *Book) DownloadAsset(src string) error {
	return DownloadAsset(src, path.Join(utils.Config.Paths.Thumbs, b.Isbn+".jpg"))
}

func (b *Book) InCollection(u *utils.User) (bool, error) {
	return db.InCollection(b.Id, u)
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
