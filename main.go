package main

import (
	"fmt"
	"io"
	"os"
)

var (
	dbURL  = os.Getenv("GOLB_DB")
	domain = os.Getenv("GOLB_DOMAIN")
)

func inject(w io.Writer, fn func(io.Writer)) func() {
	return func() { fn(w) }
}

func usage(w io.Writer) {
	fmt.Fprint(w, "Usage: golb <CMD> [OPTS]\n", os.Args[0])
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, "  CLI for hosting and posting golb")
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, "Commands:\n")
	fmt.Fprint(w, "  host			Host the golb server")
	fmt.Fprint(w, "  post			Post to golb")
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
	case "host":
		hostCmd(os.Args[1:])
	case "post":
		postCmd(os.Args[1:])
	default:
		fmt.Println("That command is not supported (--help)")
		os.Exit(2)
	}
}
