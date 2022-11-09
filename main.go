package main

import (
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
	fmt.Fprint(w, "Usage: golb <CMD> [OPTS]\n", os.Args[0])
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, "  CLI for hosting and posting golb")
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, "Commands:\n")
	fmt.Fprint(w, "  up		Host the server")
	fmt.Fprint(w, "  mk		Make a post")
	fmt.Fprint(w, "  dl		Download a post")
	fmt.Fprint(w, "  ed		Edit a post")
	fmt.Fprint(w, "  rm		Remove a post")
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, "Options:\n")
	fmt.Fprint(w, "  -h, --help		Show help")
}

func main() {
	if len(os.Args) < 2 {
		usage(os.Stderr)
		os.Exit(2)
	}

	switch os.Args[1] {
	case "up":
		upCmd(os.Args[1:])
	case "mk":
		mkCmd(os.Args[1:])
	case "dl":
		dlCmd(os.Args[1:])
	case "ed":
		edCmd(os.Args[1:])
	case "rm":
		rmCmd(os.Args[1:])
	default:
		fmt.Println("That command is not supported (--help)")
		os.Exit(2)
	}
}
