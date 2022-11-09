package main

import (
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/doug-martin/goqu/v9"
	_ "modernc.org/sqlite"
)

func dlUsage(w io.Writer) {
	fmt.Fprint(w, "Usage: golb dl <pid>\n")
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, "  Download a post")
}

func dlCmd(db *sql.DB, args []string) {
	// Parse and validate inputs
	fs := flag.NewFlagSet("dl", flag.ExitOnError)
	fs.Usage = inject(os.Stderr, edUsage)
	if err := fs.Parse(args[1:]); err != nil {
		fmt.Printf("error parsing flags (maybe consult -h): %s", err.Error())
		os.Exit(2)
	}

	pid := fs.Arg(0)
	if pid == "" {
		fmt.Println("You must supply a post ID")
		os.Exit(2)
	}

	// Get the post
	sel := goqu.Dialect("sqlite").
		Select(
			goqu.C("id"),
			goqu.C("title"),
			goqu.C("tags"),
			goqu.C("content"),
			goqu.C("created")).
		From(goqu.T("posts")).
		Where(goqu.I("id").Eq(pid))
	q, a, err := sel.ToSQL()
	if err != nil {
		fmt.Printf("error building selection query: %s\n", err.Error())
		os.Exit(1)
	}

	var post Post
	if err := db.QueryRow(q, a...).Scan(
		&post.ID,
		&post.Title,
		&post.Tags,
		&post.Content,
		&post.Created,
	); err != nil {
		fmt.Printf("error executing selection query: %s\n", err.Error())
		os.Exit(1)
	}

	// Write the post content to STDOUT
	if _, err := io.Copy(os.Stdout, bytes.NewBufferString(post.Content)); err != nil {
		fmt.Printf("error writing post content out: %s\n", err.Error())
		os.Exit(1)
	}
}
