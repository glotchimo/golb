package main

import (
	"database/sql"
	_ "embed"
	"fmt"
	"io"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	//go:embed posts.sql
	postsSQL string
	conf     Conf
)

type Conf struct {
	Name        string `yaml:"name"`
	Domain      string `yaml:"domain"`
	Port        string `yaml:"port"`
	DatabaseURL string `yaml:"database_url"`
}

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

	// Open and parse config file
	f, err := os.Open("conf.yml")
	if err != nil {
		fmt.Printf("error opening config: %s", err.Error())
		os.Exit(2)
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(&conf); err != nil {
		fmt.Printf("error parsing config: %s", err.Error())
		os.Exit(1)
	}

	// Setup SQLite DB/connection
	db, err := sql.Open("sqlite", conf.DatabaseURL)
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
