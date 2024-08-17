package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	_ = godotenv.Load()
	libsql_url := os.Getenv("libsql_url")
	libsql_token := os.Getenv("libsql_token")
	url := fmt.Sprintf("%s?authToken=%s", libsql_url, libsql_token)

	db, err := sqlx.Open("libsql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
		os.Exit(1)
	}
	defer db.Close()

	fs := http.FileServer(http.Dir("./assets"))

	state := appState{
		db: db,
	}

	http.HandleFunc("GET /", state.homeHandler)
	http.HandleFunc("GET /page/{id}", state.pageHandler)
	http.HandleFunc("GET /post/{id}", state.postHandler)
	http.Handle("GET /assets/", http.StripPrefix("/assets/", fs))

	http.ListenAndServe(":3000", nil)
}
