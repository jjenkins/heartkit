package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/zenazn/goji/web"
)

type notification struct {
	CollectionType string `json:"collectionType"`
	Date           string `json:"date"`
	OwnerId        string `json:"ownerId"`
	OwnerType      string `json:"ownerType"`
	SubscriptionId string `json:"subscriptionId"`
}

func Subscriber(c web.C, w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	var n []notification
	if err := json.Unmarshal(body, &n); err != nil {
		panic(err)
	}

	for i := 0; i < len(n); i++ {
		log.Printf("Incoming notification: %+v", n[i].SubscriptionId)
	}

	w.WriteHeader(http.StatusNoContent)
}
