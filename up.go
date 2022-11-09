package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/doug-martin/goqu/v9"
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

		// Write blog name
		if name != "" {
			var border string
			for i := 0; i < len(name); i++ {
				border += "-"
			}

			fmt.Fprintf(w, "%s |\n%s-+\n\n", name, border)
		}

		// Write posts
		for _, p := range posts {
			fmt.Fprintf(w, "%s\n", p.Title)
			fmt.Fprintf(w, "%s\n", p.Created.String())
			fmt.Fprintf(w, "%s\n", p.Tags)
			fmt.Fprintf(w, "================================\n")
			fmt.Fprintf(w, "%s\n\n\n", p.Content)
		}

		// Write footer
		fmt.Fprintf(w, "brought to you by [golb](ssh://git.glotchimo.dev/golb) with <3")
	})

	http.ListenAndServe(":8080", nil)
}
