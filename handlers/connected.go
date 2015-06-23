package handlers

import (
	"net/http"

	"github.com/zenazn/goji/web"
)

func Connected(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Connected!"))
}
