package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jjenkins/heartkit/models"
	"github.com/zenazn/goji/web"
	"golang.org/x/oauth2"
)

type notification struct {
	CollectionType string `json:"collectionType"`
	Date           string `json:"date"`
	OwnerId        string `json:"ownerId"`
	OwnerType      string `json:"ownerType"`
	SubscriptionId string `json:"subscriptionId"`
}

func Subscriber(c web.C, w http.ResponseWriter, r *http.Request) {

	apiHost := "https://api.fitbit.com"

	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	var n []notification
	if err := json.Unmarshal(body, &n); err != nil {
		panic(err)
	}

	conf := c.Env["oauth_config"].(*oauth2.Config)

	for i := 0; i < len(n); i++ {
		id := n[i].SubscriptionId
		log.Printf("Incoming notification: %+v", n[i])

		token, err := models.LoadToken(c, id)

		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			client := conf.Client(oauth2.NoContext, token)
			res, err := client.Get(fmt.Sprintf("%s/1/user/-/profile.json", apiHost))

			log.Printf("Notification response: %+v", res)
			log.Printf("Notification error: %+v", err)
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
