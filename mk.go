package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/rs/xid"
	_ "modernc.org/sqlite"
)

func mkUsage(w io.Writer) {
	fmt.Fprint(w, "Usage: golb mk [OPTS]\n")
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, "  Make a post")
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, "Options:\n")
	fmt.Fprint(w, "  -t, --title		Post title")
	fmt.Fprint(w, "  -g, --tags			Post tags")
	fmt.Fprint(w, "  -p, --path			Path to post content (.md file)")
}

func mkCmd(db *sql.DB, args []string) {
	// Parse and validate inputs
	var title, tags, path string

	fs := flag.NewFlagSet("mk", flag.ExitOnError)
	fs.Usage = inject(os.Stderr, mkUsage)

	fs.StringVar(&title, "t", "", "")
	fs.StringVar(&title, "title", "", "")

	fs.StringVar(&tags, "g", "", "")
	fs.StringVar(&tags, "tags", "", "")

	fs.StringVar(&path, "p", "", "")
	fs.StringVar(&path, "path", "", "")

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
