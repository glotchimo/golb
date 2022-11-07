package main

import (
	"database/sql"
	_ "embed"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/rs/xid"
	_ "modernc.org/sqlite"
)

//go:embed posts.sql
var postsSQL string

type Post struct {
	ID      string
	Hash    string
	Title   string
	Tags    string
	Content string
	Created time.Time
}

func postUsage(w io.Writer) {
	fmt.Fprint(w, "Usage: golb post [OPTS]\n")
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, "  Post to golb")
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, "Options:\n")
	fmt.Fprint(w, "  -t, --title		Post title")
	fmt.Fprint(w, "  -g, --tags			Post tags")
	fmt.Fprint(w, "  -p, --path			Path to post content (.md file)")
}

func postCmd(args []string) {
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
	var title, tags, path string

	fs := flag.NewFlagSet("post", flag.ExitOnError)
	fs.Usage = inject(os.Stderr, postUsage)

	fs.StringVar(&title, "t", "", "Title")
	fs.StringVar(&title, "title", "", "Title")

	fs.StringVar(&tags, "g", "", "Tags (comma-separated)")
	fs.StringVar(&tags, "tags", "", "Tags (comma-separated)")

	fs.StringVar(&path, "p", "", "Path to post content (i.e. .md file)")
	fs.StringVar(&path, "path", "", "Path to post content (i.e. .md file)")

	if err := fs.Parse(args[1:]); err != nil {
		fmt.Printf("error parsing flags (maybe consult -h): %s", err.Error())
		os.Exit(2)
	}

	if title == "" {
		fmt.Println("You must supply a title (-t)")
		os.Exit(2)
	}

	if path == "" {
		fmt.Println("You must supply a path to the post content (-p)")
		os.Exit(2)
	}

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

	// Create the post
	post := Post{
		ID:      xid.New().String(),
		Title:   title,
		Tags:    tags,
		Content: string(content),
		Created: time.Now(),
	}

	ins := goqu.Dialect("sqlite").
		Insert(goqu.T("posts")).
		OnConflict(goqu.DoNothing()).
		Rows(goqu.Record{
			"id":      post.ID,
			"title":   post.Title,
			"content": post.Content,
			"tags":    post.Tags,
			"created": post.Created,
		})
	q, a, err := ins.ToSQL()
	if err != nil {
		fmt.Printf("error building insertion query: %s\n", err.Error())
		os.Exit(1)
	}

	if _, err := db.Exec(q, a...); err != nil {
		fmt.Printf("error executing insertion query: %s\n", err.Error())
		os.Exit(1)
	}

	// Report and exit
	fmt.Printf("posted %s to %s/%s\n", path, domain, post.ID)
	os.Exit(0)
}
