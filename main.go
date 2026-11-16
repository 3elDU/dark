package main

import (
	"flag"
	"fmt"

	"github.com/3elDU/dark/browser"
	"github.com/3elDU/dark/language"
	"github.com/3elDU/dark/useragent"
)

var ua = flag.String("u", "dark", "User agent to use for all HTTP requests.\nOne of: firefox, chrome, dark")

func main() {
	flag.Parse()

	if *ua != "dark" {
		useragent.OverrideUserAgent(*ua)
	}
	fmt.Println("running with user agent:", useragent.GetUserAgent())
	fmt.Printf("running with language=%v, languages=%v\n", language.GetUserLanguage(), language.GetUserLanguages())

	url := flag.Arg(0)
	if url == "" {
		panic("expected a url as the last argument")
	}

	tab := browser.NewBrowsingContext()
	tab.NavigateTo(url)

	tab.ActiveDocument.Realm.Dom.Print()
	defer tab.Close()

	tab.ActiveDocument.Realm.MainLoop()
}
