package browser

import (
	"fmt"
	"net/url"

	"github.com/3elDU/dark/dom"
	"github.com/3elDU/dark/dom/parser"
	"github.com/3elDU/dark/http"
	"github.com/3elDU/dark/webapi"
	"github.com/grafana/sobek"
)

type Document struct {
	Realm *webapi.Realm
}

func NewBlankDocument(tab *Tab) *Document {
	blankUrl, _ := url.Parse("about:blank")

	return &Document{
		Realm: webapi.NewRealm(tab, blankUrl, dom.NewDom()),
	}
}

func NewHtmlDocument(tab *Tab, docUrl string, html string) *Document {
	doc := NewBlankDocument(tab)

	dom, err := parser.ParseDocument(html)
	if err != nil {
		panic(err)
	}

	urlParsed, err := url.Parse(docUrl)
	if err != nil {
		panic(err)
	}

	doc.Realm.Url = urlParsed
	doc.Realm.Dom = dom

	doc.loadScripts()

	return doc
}

func (d *Document) Close() {
	d.Realm.Close()
}

func (d *Document) loadScripts() {
	scripts := d.Realm.Dom.QueryAllByTagName("script")
	fmt.Printf("found %v scripts on the webpage\n", len(scripts))

	for _, script := range scripts {
		src, normUrl := d.fetchScriptSourceCode(script)

		prog, err := sobek.Compile(normUrl, src, true)
		if err != nil {
			fmt.Printf("[%v] skipping script because of compile error: %v\n", normUrl, err)
			continue
		}

		fmt.Printf("[%v] running script\n", normUrl)

		if res, err := d.Realm.VM.RunProgram(prog); err != nil {
			fmt.Printf("error when running script: %v\n", err)
		} else {
			fmt.Printf("script returned: %v\n", res.String())
		}
	}
}

func (d *Document) fetchScriptSourceCode(script *dom.Node) (source, normalizedUrl string) {
	srcUrl, hasSrc := script.Data.Attributes()["src"]

	if !hasSrc {
		return d.Realm.Dom.GetTextContent(script), "inline script"
	}

	var scriptUrl *url.URL

	if isAbsoluteUrl(srcUrl) {
		scriptUrl, _ = url.Parse(srcUrl)
	} else {
		scriptUrl = d.Realm.Url.Clone()
		scriptUrl.Path = srcUrl
	}

	source, err := http.FetchWebResource(scriptUrl.String())
	if err != nil {
		panic(err)
	}

	return source, scriptUrl.String()
}

func isAbsoluteUrl(text string) bool {
	url, err := url.Parse(text)
	if err != nil {
		return false
	}

	return url.IsAbs()
}

func createJsRealm() *sobek.Runtime {
	vm := sobek.New()

	obj := vm.NewObject()

	obj.Set("querySelector", func(call sobek.FunctionCall) sobek.Value {
		fmt.Printf("querySelector(%v)\n", call.Arguments[0].String())
		return sobek.Undefined()
	})
	obj.Set("querySelectorAll", func(call sobek.FunctionCall) sobek.Value {
		fmt.Printf("querySelectorAll(%v)\n", call.Arguments[0].String())
		return sobek.Undefined()
	})

	vm.Set("window", obj)
	vm.Set("document", obj)
	vm.Set("self", obj)

	return vm
}
