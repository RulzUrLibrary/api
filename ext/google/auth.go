package google

import (
	"github.com/rulzurlibrary/api/ext/db"
	"github.com/rulzurlibrary/api/utils"
	"net/http"
	"time"
)

const endpoint = "https://www.googleapis.com/oauth2/v3/tokeninfo"
const aud = "420546501001-v35bges8923km4s9r9p3tet8m42ibj5m.apps.googleusercontent.com"

type validationResp struct {
	AUD string `json:"aud"`
}

func Auth(db *db.DB, name, token string) (*utils.User, error) {
	var validation validationResp

	url := utils.UrlQueryS(endpoint, "id_token", token)
	resp, err := (&http.Client{Timeout: 20 * time.Second}).Get(url.String())

	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, utils.ErrGoogleAuth
	}
	if err := utils.DecodeJSON(resp.Body, &validation); err != nil {
		return nil, err
	}
	if validation.AUD != aud {
		return nil, utils.ErrGoogleAuth
	}

	return db.AuthGoogle(name)
}
