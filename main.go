package main

import (
	"blog-go/routes"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	_ = godotenv.Load()

	searchTemplate, err := template.ParseFiles("templates/base.html", "templates/search.html")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load template %v", err)
		os.Exit(1)
	}

	pageTemplate, err := template.ParseFiles("templates/base.html", "templates/page.html")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load template %v", err)
		os.Exit(1)
	}

	postTemplate, err := template.ParseFiles("templates/base.html", "templates/post.html")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load template %v", err)
		os.Exit(1)
	}

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

	state := routes.AppState{
		Db:             db,
		SearchTemplate: searchTemplate,
		PageTemplate:   pageTemplate,
		PostTemplate:   postTemplate,
	}

	http.HandleFunc("GET /", state.HomeHandler)
	http.HandleFunc("GET /page/{id}", state.PageHandler)
	http.HandleFunc("GET /post/{id}", state.PostHandler)
	http.Handle("GET /assets/", http.StripPrefix("/assets/", fs))

	http.ListenAndServe(":3000", nil)
}
