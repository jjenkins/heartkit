package handlers

import (
	"net/http"

	"github.com/zenazn/goji/web"
)

func Root(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HeartKit... powered by AliveCor"))

}
