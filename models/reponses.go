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
	Md      sql.NullString
}

type Page struct {
	Title   string
	Content template.HTML
	Url     string
}

type PostData struct {
	Slug      string
	Title     string
	Tag       sql.NullString
	Content   sql.NullString
	Md        sql.NullString
	Desc      sql.NullString
	Timestamp sql.Null[int32]
	Published bool
}

type Post struct {
	Title   string
	Tag     string
	Content template.HTML
	Url     string
	Desc    string
}

type PostRowData struct {
	Slug  string
	Title string
	Tag   string
	Desc  sql.NullString
}

type PostRow struct {
	Slug  string
	Title string
	Tag   string
	Desc  string
}

type SearchResults struct {
	Title string
	Prev  string
	Next  string
	Url   template.HTML
	Posts []PostRow
}
