package main

import (
	"database/sql"
	_ "embed"
	"fmt"
	"io"
	"os"
	"time"
)

var (
	//go:embed posts.sql
	postsSQL string

	dbURL  = os.Getenv("GOLB_DB")
	domain = os.Getenv("GOLB_DOMAIN")
	name   = os.Getenv("GOLB_NAME")
	port   = os.Getenv("GOLB_PORT")
)

type Post struct {
	ID      string
	Title   string
	Tags    string
	Content string
	Created time.Time
}

func inject(w io.Writer, fn func(io.Writer)) func() {
	return func() { fn(w) }
}

func usage(w io.Writer) {
	fmt.Fprint(w, "Usage: golb <CMD> [OPTS]\n")
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, "  CLI for hosting and managing a golb blog\n")
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, "Commands:\n")
	fmt.Fprint(w, "  up		Host the server\n")
	fmt.Fprint(w, "  mk		Make a post\n")
	fmt.Fprint(w, "  mk		List posts\n")
	fmt.Fprint(w, "  dl		Download a post\n")
	fmt.Fprint(w, "  ed		Edit a post\n")
	fmt.Fprint(w, "  rm		Remove a post\n")
}

func main() {
	if len(os.Args) < 2 {
		usage(os.Stderr)
		os.Exit(2)
	}

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

	switch os.Args[1] {
	case "up":
		upCmd(db)
	case "mk":
		mkCmd(db, os.Args[1:])
	case "ls":
		lsCmd(db)
	case "dl":
		dlCmd(db, os.Args[1:])
	case "ed":
		edCmd(db, os.Args[1:])
	case "rm":
		rmCmd(db, os.Args[1:])
	default:
		fmt.Println("That command is not supported (--help)")
		os.Exit(2)
	}
}
