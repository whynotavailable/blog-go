package models

import (
	"database/sql"
	"html/template"
)

type PageData struct {
	Id      string
	Title   string
	Excerpt sql.Null[string]
	Content sql.Null[string]
}

type Page struct {
	Title   string
	Content template.HTML
}

type PostData struct {
	Slug      string
	Title     string
	Tag       sql.Null[string]
	Content   sql.Null[string]
	Timestamp sql.Null[int32]
	Published bool
}

type Post struct {
	Title   string
	Tag     string
	Content template.HTML
}

type PostRow struct {
	Slug  string
	Title string
	Tag   string
}

type SearchResults struct {
	Title string
	Prev  string
	Next  string
	Posts []PostRow
}
