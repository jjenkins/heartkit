package handlers

import (
	"net/http"

	"github.com/zenazn/goji/web"
	"golang.org/x/oauth2"
)

func Authorize(c web.C, w http.ResponseWriter, r *http.Request) {
	conf := c.Env["oauth_config"].(*oauth2.Config)
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
