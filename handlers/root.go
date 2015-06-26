package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jjenkins/heartkit/models"
	"github.com/zenazn/goji/web"
	"golang.org/x/oauth2"
)

func Root(c web.C, w http.ResponseWriter, r *http.Request) {

	apiHost := "https://api.fitbit.com"
	conf := c.Env["oauth_config"].(*oauth2.Config)
	id := "lfws3agqhra2vfkpjlb4e6sxie"
	token, _ := models.LoadToken(c, id)

	url := fmt.Sprintf("%s/1/user/-/activities/apiSubscriptions/%s-activities.json", apiHost, id)
	deleteRequest, err := http.NewRequest("DELETE", url, nil)

	if err != nil {
		log.Printf("Error: %v", err)
	}

	client := conf.Client(oauth2.NoContext, token)
	res, err := client.Do(deleteRequest)

	if err != nil {
		log.Printf("Error: %v", err)
	}

	log.Printf("Response: %v", res)

	w.Write([]byte("HeartKit... powered by AliveCor"))

}
