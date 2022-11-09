package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/gomarkdown/markdown"
)

func upUsage(w io.Writer) {
	fmt.Fprint(w, "Usage: golb up\n")
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, "  Host the server")
}

func upCmd(db *sql.DB) {
	// Handle requests by returning post contents in a list
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Get the posts
		sel := goqu.Dialect("sqlite").
			Select(
				goqu.C("id"),
				goqu.C("title"),
				goqu.C("tags"),
				goqu.C("content"),
				goqu.C("created")).
			From(goqu.T("posts")).
			Order(goqu.I("id").Desc())
		q, a, err := sel.ToSQL()
		if err != nil {
			fmt.Printf("error building selection query: %s\n", err.Error())
			os.Exit(1)
		}

		var posts []Post
		rows, err := db.Query(q, a...)
		if err != nil {
			fmt.Printf("error executing selection query: %s\n", err.Error())
			os.Exit(1)
		}

		for rows.Next() {
			var post Post
			if err := rows.Scan(
				&post.ID,
				&post.Title,
				&post.Tags,
				&post.Content,
				&post.Created,
			); err != nil {
				fmt.Printf("error scanning post: %s\n", err.Error())
				os.Exit(1)
			}

			posts = append(posts, post)
		}

		if err := rows.Err(); err != nil {
			fmt.Printf("error when processing rows: %s\n", err.Error())
			os.Exit(1)
		}

		// Combine page content
		var page string

		if name != "" {
			page += "# " + name + "\n\n---\n\n"
		}

		for _, p := range posts {
			page += "## " + p.Title + "\n\n"
			page += "*" + p.Created.Format(time.RFC1123) + "* // *" + p.Tags + "*" + "\n\n"
			page += p.Content
			page += "\n\n---\n\n"
		}
		page += "brought to you by [golb](ssh://git.glotchimo.dev/golb) with <3"

		md := []byte(page)
		output := markdown.ToHTML(md, nil, nil)
		w.Write(output)
	})

	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, nil)
}
