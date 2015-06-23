package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/zenazn/goji/web"
)

func Subscriber(c web.C, w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	var notification map[string]interface{}
	if err := json.Unmarshal(body, &notification); err != nil {
		panic(err)
	}

	log.Printf("Incoming notification: %+v", notification)
	w.WriteHeader(http.StatusNoContent)
}
