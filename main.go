package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/alivecor/surge/handlers"
	_ "github.com/lib/pq"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"golang.org/x/oauth2"
)

func config(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		conf := &oauth2.Config{
			ClientID:     os.Getenv("FITBIT_CLIENT_ID"),
			ClientSecret: os.Getenv("FITBIT_CLIENT_SECRET"),
			Scopes:       []string{"activity", "heartrate"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://www.fitbit.com/oauth2/authorize",
				TokenURL: "https://api.fitbit.com/oauth2/token",
			},
		}

		c.Env["oauth_config"] = conf
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func database(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		db, _ := sql.Open("postgres", os.Getenv("DATABASE_URL"))
		err := db.Ping()

		if err != nil {
			log.Fatal(err)
		}

		c.Env["db"] = db
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func main() {
	goji.Use(config)
	goji.Use(database)
	goji.Get("/", handlers.Root)
	goji.Get("/authorize", handlers.Authorize)
	goji.Get("/connected", handlers.Connected)
	goji.Get("/callback", handlers.Callback)
	goji.Post("/subscriber", handlers.Subscriber)
	goji.NotFound(handlers.NotFound)
	goji.Serve()
}
