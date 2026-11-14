package browser

import (
	"github.com/3elDU/dark/dom"
	"github.com/3elDU/dark/dom/parser"
	"github.com/grafana/sobek"
)

type Document struct {
	Dom *dom.Dom
	Js  *sobek.Runtime
}

func NewBlankDocument() *Document {
	return &Document{
		Dom: dom.NewDom(),
		Js:  createJsRealm(),
	}
}

func NewHtmlDocument(html string) *Document {
	dom, err := parser.ParseDocument(html)
	if err != nil {
		panic(err)
	}

	return &Document{
		Dom: dom,
		Js:  createJsRealm(),
	}
}

func createJsRealm() *sobek.Runtime {
	return sobek.New()
}
