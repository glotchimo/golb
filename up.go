package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/doug-martin/goqu/v9"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

func upUsage(w io.Writer) {
	fmt.Fprint(w, "Usage: golb up\n")
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, "  Host the server")
}

func upCmd(db *sql.DB) {
	// Handle requests by returning post contents in a list
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s: %s %s (%dB)", r.RemoteAddr, r.Method, r.Host, r.ContentLength)

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

		page := buildContent(conf, posts)

		// Compile markdown to HTML and write out
		md := []byte(page)
		extensions := parser.CommonExtensions | parser.AutoHeadingIDs
		parser := parser.NewWithExtensions(extensions)
		output := markdown.ToHTML(md, parser, nil)
		w.Write(output)
	})

	if conf.Port == "" {
		conf.Port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+conf.Port, nil))
}
