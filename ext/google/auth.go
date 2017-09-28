package google

import (
	"github.com/ixday/echo-hello/ext/db"
	"github.com/ixday/echo-hello/utils"
	"net/http"
	"time"
)

const endpoint = "https://www.googleapis.com/oauth2/v3/tokeninfo"
const aud = "420546501001-v35bges8923km4s9r9p3tet8m42ibj5m.apps.googleusercontent.com"

type validationResp struct {
	AUD string `json:"aud"`
}

func Auth(db *db.DB, name, token string) (u utils.User, _ error) {
	var validation validationResp

	url := utils.UrlQueryS(endpoint, "id_token", token)
	resp, err := (&http.Client{Timeout: 20 * time.Second}).Get(url.String())

	if err != nil {
		return u, err
	}
	if resp.StatusCode != http.StatusOK {
		return u, utils.ErrGoogleAuth
	}
	if err := utils.DecodeJSON(resp.Body, &validation); err != nil {
		return u, err
	}
	if validation.AUD != aud {
		return u, utils.ErrGoogleAuth
	}

	return db.AuthGoogle(name)
}
