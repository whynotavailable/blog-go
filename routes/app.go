package routes

import (
	"blog-go/models"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Data interface {
	Select(dest interface{}, query string, args ...interface{}) error
	Get(dest interface{}, query string, args ...interface{}) error
}

type AppState struct {
	Db             Data
	SearchTemplate *template.Template
	PageTemplate   *template.Template
	PostTemplate   *template.Template
}

func getPrev(page int, tag string) string {
	result := ""
	if page > 0 {
		if page == 1 {
			result = "/"
		} else if page > 1 {
			result = fmt.Sprintf("/?page=%d", page-1)
		}

		if tag != "" {
			if strings.Contains(result, "?") {
				result += "&"
			} else {
				result += "?"
			}

			result += fmt.Sprintf("tag=%s", url.QueryEscape(tag))
		}
	}

	return result
}

func (state *AppState) HomeHandler(w http.ResponseWriter, r *http.Request) {
	con_url := fmt.Sprintf("https://whynot.sh%s", r.URL.String())
	fmt.Println(con_url)

	query := r.URL.Query()

	page, err := strconv.Atoi(query.Get("page"))
	if err != nil {
		page = 0
	}

	tag := query.Get("tag")

	posts := []models.PostRow{}
	searchResults := models.SearchResults{
		Url: template.HTML(con_url),
	}

	offset := page * 5

	if tag == "" {
		// SELECT slug, tag, title FROM post WHERE published = TRUE AND tag = ?1 ORDER BY timestamp DESC LIMIT 6 OFFSET ?2
		state.Db.Select(&posts, "SELECT slug, tag, title FROM post WHERE published = TRUE ORDER BY timestamp DESC LIMIT 6 OFFSET ?1", offset)
		searchResults.Title = "Home"
	} else {
		state.Db.Select(&posts, "SELECT slug, tag, title FROM post WHERE published = TRUE AND tag = ?1 ORDER BY timestamp DESC LIMIT 6 OFFSET ?2", tag, offset)
		searchResults.Title = fmt.Sprintf("%s List", tag)
	}

	// Previous Page
	searchResults.Prev = getPrev(page, tag)

	// Next Page
	if len(posts) > 5 {
		searchResults.Next = fmt.Sprintf("/?page=%d", page+1)

		if tag != "" {
			searchResults.Next += fmt.Sprintf("&tag=%s", url.QueryEscape(tag))
		}

		posts = posts[:5]
	}

	searchResults.Posts = posts

	err = state.SearchTemplate.Execute(w, searchResults)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "%v", err)
		return
	}
}

func (state *AppState) PageHandler(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("https://whynot.sh%s", r.URL.String())

	pageModel := models.PageData{}
	err := state.Db.Get(&pageModel, "SELECT * FROM page WHERE id = ?1", r.PathValue("id"))
	if err != nil {
		// TODO: Figure out how to make this 404 properly
		w.WriteHeader(500)
		fmt.Fprintf(w, "%v", err)
		return
	}

	page := models.Page{
		Title: pageModel.Title,
		Url:   url,
	}

	if pageModel.Content.Valid {
		page.Content = template.HTML(pageModel.Content.V)
	}

	// hello := templates.Hello(r.PathValue("id"))
	err = state.PageTemplate.Execute(w, page)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "%v", err)
		return
	}
}

func (state *AppState) PostHandler(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("https://whynot.sh%s", r.URL.String())

	pageModel := models.PostData{}
	err := state.Db.Get(&pageModel, "SELECT * FROM post WHERE slug = ?1", r.PathValue("id"))
	if err != nil {
		// TODO: Figure out how to make this 404 properly
		w.WriteHeader(500)
		fmt.Fprintf(w, "%v", err)
		return
	}

	page := models.Post{
		Title: pageModel.Title,
		Url:   url,
	}

	if pageModel.Tag.Valid {
		page.Tag = pageModel.Tag.String
	} else {
		page.Tag = "tbd"
	}

	if pageModel.Desc.Valid {
		page.Desc = pageModel.Desc.String
	}

	if pageModel.Content.Valid {
		page.Content = template.HTML(pageModel.Content.String)
	}

	// hello := templates.Hello(r.PathValue("id"))
	err = state.PostTemplate.Execute(w, page)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "%v", err)
		return
	}
}
