package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Sam-Izdat/kee"
	"github.com/jjenkins/heartkit/models"
	"github.com/zenazn/goji/web"
	"golang.org/x/oauth2"
)

func Callback(c web.C, w http.ResponseWriter, r *http.Request) {

	apiHost := "https://api.fitbit.com"

	conf := c.Env["oauth_config"].(*oauth2.Config)
	code := r.URL.Query().Get("code")
	tok, err := conf.Exchange(oauth2.NoContext, code)

	if err != nil {
		log.Fatal(err)
	}

	client := conf.Client(oauth2.NoContext, tok)
	res, err := client.Get(fmt.Sprintf("%s/1/user/-/profile.json", apiHost))

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	type user struct {
		Profile models.User `json:"user"`
	}

	u := &user{}
	if err := json.Unmarshal(body, u); err != nil {
		panic(err)
	}

	// Look up the user and see if we've seen this person
	// before. If so, update the row, otherwise insert a
	// new user into the table.
	db := c.Env["db"].(*sql.DB)
	tx, _ := db.Begin()

	query, err := db.Prepare("SELECT id FROM users WHERE encoded_id = $1 LIMIT 1")

	var encodedId string
	err = query.QueryRow(u.Profile.EncodedId).Scan(&encodedId)

	if err == sql.ErrNoRows {
		// Insert user into database.
		query := `
			INSERT INTO users
				(id, encoded_id, gender, date_of_birth, access_token,
					refresh_token, expiry, token_type, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`

		stmt, err := tx.Prepare(query)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer stmt.Close()

		kee.UUID.Options.HyphURL32 = false
		id := kee.UUID.New()

		_, err = stmt.Exec(
			strings.ToLower(id.URL32()),
			u.Profile.EncodedId,
			u.Profile.Gender,
			u.Profile.DateOfBirth,
			tok.AccessToken,
			tok.RefreshToken,
			tok.Expiry,
			tok.TokenType,
			time.Now(),
			time.Now(),
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Setup the subscriptions
		url := fmt.Sprintf("%s/1/user/%s/activities/apiSubscriptions/%s.json",
			apiHost, u.Profile.EncodedId, strings.ToLower(id.URL32()))

		res, err := client.Post(url, "", nil)
		log.Printf("Activities subscription request: %+v", res)

		url = fmt.Sprintf("%s/1/user/%s/body/apiSubscriptions/%s.json",
			apiHost, u.Profile.EncodedId, strings.ToLower(id.URL32()))

		res, err = client.Post(url, "", nil)
		log.Printf("Body subscription request: %+v", res)

	} else if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	} else {

		query := `
			UPDATE users
			SET
				gender=$1,
				date_of_birth=$2,
				access_token=$3,
				refresh_token=$4,
				expiry=$5,
				token_type=$6,
				updated_at=$7
			WHERE encoded_id=$8
		`
		stmt, err := tx.Prepare(query)
		defer stmt.Close()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = stmt.Exec(
			u.Profile.Gender,
			u.Profile.DateOfBirth,
			tok.AccessToken,
			tok.RefreshToken,
			tok.Expiry,
			tok.TokenType,
			time.Now(),
			u.Profile.EncodedId,
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/connected", 301)
}
