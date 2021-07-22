package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/gorilla/pat"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	io.WriteString(w, `{"alive": true}`)
}

func TestUser(res http.ResponseWriter, req *http.Request) {
	user, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		fmt.Fprintln(res, err)
		return
	}
	foo, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Error", err)
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(foo)
	http.Redirect(res, req, "http://localhost:3000/auth/show", 302)
}

func AuthUser(res http.ResponseWriter, req *http.Request) {
	if gothUser, err := gothic.CompleteUserAuth(res, req); err == nil {
		// fmt.Println("yha aya", gothUser)
		res.Header().Set("Content-Type", "application/json")
		_, err := json.Marshal(gothUser)
		if err != nil {
			fmt.Println("Error", err)
		}

		// fmt.Println(gothUser)
		json.NewEncoder(res).Encode(gothUser)
		// res.Write(foo)
		http.Redirect(res, req, "http://localhost:3000/auth/show", 302)
	} else {
		gothic.BeginAuthHandler(res, req)
	}
	// http.Redirect(res, req, "http://localhost:3000/auth/show", 302)
}

func main() {
	goth.UseProviders(
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), "http://localhost:4200/auth/github/callback"),
	)
	p := pat.New()
	p.Get("/health", HealthCheckHandler)
	p.Get("/auth/{provider}", AuthUser)
	// p.Post("/auth/{provider}/callback", TestUser)
	router := mux.NewRouter()
	router.HandleFunc("/auth/{provider}/callback", TestUser).Methods("POST")

	log.Println("listening on localhost:4200")
	log.Fatal(http.ListenAndServe("localhost:4200", p))
}
