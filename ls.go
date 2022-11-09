package main

import (
	"database/sql"
	"fmt"
	"io"
	"os"

	"github.com/doug-martin/goqu/v9"
	_ "modernc.org/sqlite"
)

func lsUsage(w io.Writer) {
	fmt.Fprint(w, "Usage: golb ls\n")
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, "  List posts")
}

func lsCmd(db *sql.DB) {
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

	// List out ID, created datetime, title, and tags
	fmt.Println("id\tcreated\ttitle\ttags")
	for _, p := range posts {
		fmt.Printf("%s\t%s\t%s\t%s\n", p.ID, p.Created.String(), p.Title, p.Tags)
	}
}
