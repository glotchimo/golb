package main

import (
	"database/sql"
	_ "embed"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/doug-martin/goqu/v9"
	_ "modernc.org/sqlite"
)

func rmUsage(w io.Writer) {
	fmt.Fprint(w, "Usage: golb rm [OPTS]\n")
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, "  Remove a post")
}

func rmCmd(db *sql.DB, args []string) {
	// Parse and validate inputs
	fs := flag.NewFlagSet("rm", flag.ExitOnError)
	fs.Usage = inject(os.Stderr, rmUsage)
	if err := fs.Parse(args[1:]); err != nil {
		fmt.Printf("error parsing flags (maybe consult -h): %s", err.Error())
		os.Exit(2)
	}

	pid := fs.Arg(0)
	if pid == "" {
		fmt.Println("You must supply a post ID")
		os.Exit(2)
	}

	// Delete the post
	del := goqu.Dialect("sqlite").
		Delete(goqu.T("posts")).
		Where(goqu.I("id").Eq(pid))
	q, a, err := del.ToSQL()
	if err != nil {
		fmt.Printf("error building deletion query: %s\n", err.Error())
		os.Exit(1)
	}

	if _, err := db.Exec(q, a...); err != nil {
		fmt.Printf("error executing deletion query: %s\n", err.Error())
		os.Exit(1)
	}

	// Report and exit
	fmt.Printf("deleted %s\n", pid)
	os.Exit(0)
}
