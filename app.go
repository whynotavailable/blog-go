package main

import (
	"blog-go/models"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type appState struct {
	db             *sqlx.DB
	searchTemplate *template.Template
	pageTemplate   *template.Template
	postTemplate   *template.Template
}

func (state *appState) homeHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	page, err := strconv.Atoi(query.Get("page"))
	if err != nil {
		page = 0
	}

	tag := query.Get("tag")

	posts := []models.PostRow{}
	searchResults := models.SearchResults{}

	offset := page * 5

	if tag == "" {
		// SELECT slug, tag, title FROM post WHERE published = TRUE AND tag = ?1 ORDER BY timestamp DESC LIMIT 6 OFFSET ?2
		state.db.Select(&posts, "SELECT slug, tag, title FROM post WHERE published = TRUE ORDER BY timestamp DESC LIMIT 6 OFFSET ?1", offset)
		searchResults.Title = "Home"
	} else {
		state.db.Select(&posts, "SELECT slug, tag, title FROM post WHERE published = TRUE AND tag = ?1 ORDER BY timestamp DESC LIMIT 6 OFFSET ?2", tag, offset)
		searchResults.Title = fmt.Sprintf("%s List", tag)
	}

	if page > 0 {
		// Previous Page
		if page == 1 {
			searchResults.Prev = "/"
		} else if page > 1 {
			searchResults.Prev = fmt.Sprintf("/?page=%d", page-1)
		}

		if tag != "" {
			if strings.Contains(searchResults.Prev, "?") {
				searchResults.Prev += "&"
			} else {
				searchResults.Prev += "?"
			}

			searchResults.Prev += fmt.Sprintf("tag=%s", url.QueryEscape(tag))
		}
	}

	// Next Page
	if len(posts) > 5 {
		searchResults.Next = fmt.Sprintf("/?page=%d", page+1)

		if tag != "" {
			searchResults.Next += fmt.Sprintf("&tag=%s", url.QueryEscape(tag))
		}

		posts = posts[:5]
	}

	fmt.Println(searchResults.Next)
	searchResults.Posts = posts

	err = state.searchTemplate.Execute(w, searchResults)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "%v", err)
		return
	}
}

func (state *appState) pageHandler(w http.ResponseWriter, r *http.Request) {
	pageModel := models.PageData{}
	err := state.db.Get(&pageModel, "SELECT * FROM page WHERE id = ?1", r.PathValue("id"))
	if err != nil {
		// TODO: Figure out how to make this 404 properly
		w.WriteHeader(500)
		fmt.Fprintf(w, "%v", err)
		return
	}

	page := models.Page{
		Title: pageModel.Title,
	}

	if pageModel.Content.Valid {
		page.Content = template.HTML(pageModel.Content.V)
	}

	// hello := templates.Hello(r.PathValue("id"))
	err = state.pageTemplate.Execute(w, page)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "%v", err)
		return
	}
}

func (state *appState) postHandler(w http.ResponseWriter, r *http.Request) {
	pageModel := models.PostData{}
	err := state.db.Get(&pageModel, "SELECT * FROM post WHERE slug = ?1", r.PathValue("id"))
	if err != nil {
		// TODO: Figure out how to make this 404 properly
		w.WriteHeader(500)
		fmt.Fprintf(w, "%v", err)
		return
	}

	page := models.Post{
		Title: pageModel.Title,
	}

	if pageModel.Tag.Valid {
		page.Tag = pageModel.Tag.V
	} else {
		page.Tag = "tbd"
	}

	if pageModel.Content.Valid {
		page.Content = template.HTML(pageModel.Content.V)
	}

	// hello := templates.Hello(r.PathValue("id"))
	err = state.postTemplate.Execute(w, page)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "%v", err)
		return
	}
}
