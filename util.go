package main

import "time"

func buildContent(conf Conf, posts []Post) string {
	// Compose page content
	var page string
	page += "# " + conf.Heading.Title + "\n"
	page += "#### " + conf.Heading.Email + " // " + conf.Heading.Git
	page += "\n\n---\n\n"
	for _, p := range posts {
		page += "## " + p.Title + "\n\n"
		page += "*" + p.Created.Format(time.RFC1123) + "* // *" + p.Tags + "*" + "\n\n"
		page += p.Content
		page += "\n\n---\n\n"
	}
	page += "brought to you by [golb](ssh://git.plain.technology:23231/golb) // meant to be viewed in a [reader](https://www.maketecheasier.com/enable-browser-reader-mode/)"

	return page
}
