# golb

A blog, but we're going backwards!

- Self-hostable
- Single binary
- Markdown content
- Managed via the command line
- SQLite

```
go install git.glotchimo.dev/golb
```

## What?

There are a lot of really awesome, fully-featured blog platforms
out there, many of which enable creators to distribute content
and monetize it really easily, but because I can and want to, I
decided to write my own (yet another) blog engine that's as simple
and dumb as possible.

Golb is a SQLite-using, Markdown-rendering, blog engine written in Go and 
managed via a CLI that allows you to write content and then post it, and 
that's about it.

## How?

Here's a little quickstart to get things going:

```
$ echo "# My First Post" > my_first_post.md
$ golb mk -t "My First Post" -g random,writing -p my_first_post.md
posted my_first_post.md to localhost/cdlhlm28ra5ccqkpjddg
$ golb up
```

Run `golb` for this dialogue:

```
Usage: golb <CMD> [OPTS]

  CLI for hosting and managing a golb blog

Commands:
  up		Host the server
  mk		Make a post
  mk		List posts
  dl		Download a post
  ed		Edit a post
  rm		Remove a post
```

Each subcommand has its own help dialogue that you can see by adding `-h`.

## Todo

- [x] Add an rm command for content removal
- [x] Add a dl command for downloading a post into a file
- [x] Add an ed command for editing posts
- [x] Add an ls command for listing all posts
- [x] Actually host something
- [x] Render Markdown into HTML
- [x] Add logging
- [ ] Add header links for shareability
- [ ] Make an FS-backed version?
