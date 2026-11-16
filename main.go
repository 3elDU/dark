package main

import (
	"flag"

	"github.com/3elDU/dark/browser"
)

var ua = flag.String("u", "firefox", "User agent to use for all HTTP requests.\nOne of: firefox, chrome, custom")

func main() {
	flag.Parse()

	url := flag.Arg(0)
	if url == "" {
		panic("expected a url as the last argument")
	}

	tab := browser.NewBrowsingContext(*ua)
	tab.NavigateTo(url)

	tab.ActiveDocument.Realm.Dom.Print()
	defer tab.Close()

	tab.ActiveDocument.Realm.MainLoop()
}
