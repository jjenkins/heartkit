package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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

	db := c.Env["db"].(*sql.DB)
	conf := c.Env["oauth_config"].(*oauth2.Config)

	for i := 0; i < len(n); i++ {

		stmt, _ := db.Prepare("SELECT oauth_access_token, oauth_refresh_token FROM users WHERE id = $1")

		var accessToken string
		var refreshToken string
		// err = stmt.QueryRow(n[i].SubscriptionId).Scan(&accessToken, &refreshToken)
		err = stmt.QueryRow("4v3gl7mgszh55cwrxvkfe2z3va").Scan(&accessToken, &refreshToken)

		tok := oauth2.Token{
			TokenType:    "Bearer",
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		log.Printf("OAuth Token: %+v", tok)

		client := conf.Client(oauth2.NoContext, &tok)
		res, err := client.Get(fmt.Sprintf("%s/1/user/-/profile.json", apiHost))

		log.Printf("Incoming notification: %+v", n[i].SubscriptionId)
		log.Printf("Notification response: %+v", res)
		log.Printf("Notification error: %+v", err)
	}

	w.WriteHeader(http.StatusNoContent)
}
