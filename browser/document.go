package browser

import (
	"fmt"
	"net/url"

	"github.com/3elDU/dark/dom"
	"github.com/3elDU/dark/dom/parser"
	"github.com/3elDU/dark/webapi"
	"github.com/grafana/sobek"
)

type Document struct {
	Realm *webapi.Realm

	scripts []scriptRecord
	// cache normalized urls to modules
	modules map[string]sobek.ModuleRecord
	// cache linked modules to urls.
	// used to preserve the base url when a script imports another script.
	reverseModules map[any]*url.URL
}

func NewBlankDocument(tab *Tab) *Document {
	blankUrl, _ := url.Parse("about:blank")

	doc := &Document{
		scripts:        []scriptRecord{},
		modules:        map[string]sobek.ModuleRecord{},
		reverseModules: map[any]*url.URL{},
	}

	doc.Realm = webapi.NewRealm(tab, doc, blankUrl, dom.NewDom())

	return doc
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

	doc.findScripts()
	doc.loadScripts()

	return doc
}

func (d *Document) Close() {
	d.Realm.Close()
}

func (d *Document) findScripts() {
	scriptElements := d.Realm.Dom.QueryAllByTagName("script")

	for _, scriptNode := range scriptElements {
		scriptType := classifyScript(scriptNode)

		switch scriptType {
		case classicScript, moduleScript:
			d.scripts = append(d.scripts, scriptRecord{
				node: scriptNode,
				t:    scriptType,
				// FIXME: all scripts are deferred right now, because we have no
				// way of running the scripts as the dom gets built. Adding non-deferred
				// scripts requires us to switch to a different HTML parser or
				// fork x/net/html
				deferred: true,
			})
		default:
			fmt.Printf("ignoring script with type %v\n", scriptNode.Data.Attributes().Get("type"))
		}
	}

	fmt.Printf("found %v <script> elements on the webpage, loading %v of them\n",
		len(scriptElements), len(d.scripts),
	)
}

func (d *Document) loadScripts() {
	for i := range d.scripts {
		script := &d.scripts[i]
		src, normUrl, err := d.fetchScriptSourceCode(*script)

		normUrlDebug := fmt.Sprintf("inline-script-%v", i)
		if normUrl != nil {
			normUrlDebug = normUrl.String()
		}

		script.name = normUrlDebug

		if err != nil {
			fmt.Printf("[%v] error while fetching script: %v", normUrlDebug, err)
			continue
		}

		script.url = normUrl

		switch script.t {
		case classicScript:
			// classic scripts are always compiled in non-strict mode
			prog, err := sobek.Compile(script.name, src, false)
			if err != nil {
				fmt.Printf("[%v] skipping classis script because of compile error: %v\n", script.name, err)
				continue
			}

			fmt.Printf("[%v] running script\n", script.name)

			if res, err := d.Realm.VM.RunProgram(prog); err != nil {
				fmt.Printf("error when running script: %v\n", err)
			} else {
				fmt.Printf("[%v] script returned: %v\n", script.name, res.String())
				script.ran = true
			}
		case moduleScript:
			p, err := d.runModuleScriptSource(*script, src)
			if err != nil {
				fmt.Printf("[%v] error when running module script: %v\n", script.name, err)
				continue
			} else {
				fmt.Printf("[%v] completed module script, promise result = %v\n", script.name, p.Result().String())
				script.ran = true
			}
		}

	}
}
