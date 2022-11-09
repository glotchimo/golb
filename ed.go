package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/doug-martin/goqu/v9"
	_ "modernc.org/sqlite"
)

func edUsage(w io.Writer) {
	fmt.Fprint(w, "Usage: golb ed <pid> <path>\n")
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, "  Edit a post")
}

func edCmd(args []string) {
	// Setup SQLite DB/connection
	db, err := sql.Open("sqlite", os.Getenv("GOLB_DB"))
	if err != nil {
		fmt.Printf("error connecting to database: %s", err.Error())
		os.Exit(1)
	}
	defer db.Close()

	if _, err := db.Exec(postsSQL); err != nil {
		fmt.Printf("error creating posts table: %s", err.Error())
		os.Exit(1)
	}

	// Parse and validate inputs
	fs := flag.NewFlagSet("ed", flag.ExitOnError)
	fs.Usage = inject(os.Stderr, edUsage)
	if err := fs.Parse(args[1:]); err != nil {
		fmt.Printf("error parsing flags (maybe consult -h): %s", err.Error())
		os.Exit(2)
	}

	if len(fs.Args()) != 2 {
		fmt.Println("You must supply a post ID and path to the post content")
		os.Exit(2)
	}

	pid := fs.Arg(0)
	path := fs.Arg(1)

	// Load and process post content
	f, err := os.Open(path)
	if err != nil {
		fmt.Printf("error opening post file: %s\n", err.Error())
		os.Exit(1)
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		fmt.Printf("error reading post file: %s\n", err.Error())
		os.Exit(1)
	}

	// Update the post
	upd := goqu.Dialect("sqlite").
		Update(goqu.T("posts")).
		Set(goqu.Record{"content": content})
	q, a, err := upd.ToSQL()
	if err != nil {
		fmt.Printf("error building update query: %s\n", err.Error())
		os.Exit(1)
	}

	if _, err := db.Exec(q, a...); err != nil {
		fmt.Printf("error executing update query: %s\n", err.Error())
		os.Exit(1)
	}

	// Report and exit
	fmt.Printf("updated %s with %s\n", pid, path)
	os.Exit(0)
}
