package models

import (
	"database/sql"
	"time"

	"github.com/zenazn/goji/web"
	"golang.org/x/oauth2"
)

func LoadToken(c web.C, id string) (*oauth2.Token, error) {

	db := c.Env["db"].(*sql.DB)
	conf := c.Env["oauth_config"].(*oauth2.Config)

	var accessToken string
	var refreshToken string
	var expiry string
	var tokenType string

	query := `
			SELECT access_token, refresh_token, expiry, token_type
			FROM users
			WHERE id = $1
		`
	stmt, _ := db.Prepare(query)
	defer stmt.Close()

	err := stmt.QueryRow(id).Scan(&accessToken, &refreshToken, &expiry, &tokenType)
	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		return nil, err
	}

	parsedExpiry, _ := time.Parse(time.RFC3339Nano, expiry)

	token := &oauth2.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Expiry:       parsedExpiry,
		TokenType:    tokenType,
	}

	source := conf.TokenSource(oauth2.NoContext, token)
	t, err := source.Token()

	if refreshToken != t.RefreshToken {

		query := `
			UPDATE users
			SET
				access_token=$1,
				refresh_token=$2,
				expiry=$3,
				token_type=$4,
				updated_at=$5
			WHERE id=$6
		`
		stmt, _ := db.Prepare(query)
		defer stmt.Close()

		_, err = stmt.Exec(
			t.AccessToken,
			t.RefreshToken,
			t.Expiry,
			t.TokenType,
			time.Now(),
			id,
		)

		if err != nil {
			return nil, err
		}

	}

	return t, err
}
