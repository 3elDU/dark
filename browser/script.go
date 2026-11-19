package browser

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/3elDU/dark/dom"
	"github.com/3elDU/dark/http"
	"github.com/3elDU/dark/mime"
	"github.com/grafana/sobek"
)

type scriptType int

const (
	// regular JS, no modules
	classicScript scriptType = iota
	// JavaScript moduleScript
	moduleScript
	// JSON import map for JS modules
	importMapScript
	// not JavaScript
	ignoreScript
)

type scriptRecord struct {
	element dom.Element
	// used for debugging
	name string
	// normalized, always absolute. nil for inline scripts
	url      *url.URL
	t        scriptType
	deferred bool
	ran      bool
}

func classifyScript(script dom.Element) scriptType {
	scriptType := strings.ToLower(
		strings.TrimSpace(
			script.Attributes().Get("type"),
		),
	)

	nomodule := script.Attributes().Has("nomodule")
	if nomodule {
		// nomodule is for older browsers to load a fallback script.
		// we support modules, so any script with 'nomodule' is ignored.
		return ignoreScript
	}

	switch {
	case mime.IsJavascriptContentType(scriptType), scriptType == "":
		return classicScript
	case scriptType == "module":
		return moduleScript
	case scriptType == "importmap":
		return importMapScript
	default:
		return ignoreScript
	}
}

// with inline scripts, returns the source code and a nil url
// with an external script (with an 'src' attribute) returns the empty string as
// source, and the normalized script url
func (d *Document) fetchScriptSourceCode(script scriptRecord) (source string, normalizedUrl *url.URL, err error) {
	srcUrl, hasSrc := script.element.Attributes()["src"]

	if !hasSrc {
		return script.element.GetTextContent(), nil, nil
	}

	scriptUrl, err := getNormalizedScriptUrl(d.Realm.Url, srcUrl)
	if err != nil {
		return "", nil, fmt.Errorf("invalid url for 'src': %v", err.Error())
	}

	var src string

	// module scripts will be fetched by the module loader
	if script.t == classicScript {
		src, err = http.FetchWebResource(scriptUrl.String())
		if err != nil {
			return "", nil, err
		}
	}

	return src, scriptUrl, nil
}

// resolves relative paths if needed
func getNormalizedScriptUrl(base *url.URL, scriptUrl string) (*url.URL, error) {
	norm, err := url.Parse(scriptUrl)
	if err != nil {
		return nil, err
	}

	if !norm.IsAbs() {
		norm = base.ResolveReference(norm)
	}

	return norm, nil
}

func (d *Document) GetModuleReferrerUrl(referrer any) *url.URL {
	if u, ok := d.reverseModules[referrer]; ok {
		return u
	}
	return d.Realm.Url
}

// returns nil if the import map doesn't match the specifier
func (d *Document) checkImportMap(specifier string) *url.URL {
	for _, imap := range d.importMaps {
		for k, v := range imap.Imports {
			if k == specifier {
				return v
			} else if strings.HasSuffix(k, "/") && strings.HasPrefix(specifier, k) {
				u, err := url.Parse(
					strings.TrimPrefix(specifier, k),
				)
				if err == nil {
					return v.ResolveReference(u)
				}
			}
		}
	}

	return nil
}

// Given a specifier like `./a.js` or simply `react`, resolves it to an absolute url.
// Error is returned for invalid import map entries and bare specifiers.
func (d *Document) ResolveModuleSpecifier(referrer any, specifier string) (u *url.URL, err error) {
	// normalize the specifier first (resolve relative paths to get an absolute url)
	if strings.HasPrefix(specifier, "/") || strings.HasPrefix(specifier, "./") || strings.HasPrefix(specifier, "../") {
		u, err = getNormalizedScriptUrl(d.GetModuleReferrerUrl(referrer), specifier)
		if err != nil {
			return
		}
	}

	// if specifier doesn't start as a path, parse it as a regular url
	u, err = url.Parse(specifier)

	// look up the import map
	key := specifier
	if err == nil {
		key = u.String()
	}

	uMap := d.checkImportMap(key)
	if uMap != nil {
		return uMap, nil
	}

	if u == nil || !u.IsAbs() {
		return nil, fmt.Errorf("unable to resolve bare specifier: %v", specifier)
	}

	return u, nil
}

func (d *Document) ResolveImportedModule(referencingScriptOrModule interface{}, specifier string) (sobek.ModuleRecord, error) {
	norm, err := d.ResolveModuleSpecifier(referencingScriptOrModule, specifier)
	if err != nil {
		return nil, err
	}

	return d.fetchModule(norm)
}

func (d *Document) fetchModule(u *url.URL) (sobek.ModuleRecord, error) {
	if cached, ok := d.modules[u.String()]; ok {
		return cached, nil
	}

	fmt.Printf("fetchModule(url='%v')\n", u)

	src, err := http.FetchWebResource(u.String())
	if err != nil {
		return nil, err
	}

	mod, err := sobek.ParseModule(u.String(), src, d.ResolveImportedModule)
	if err == nil {
		d.modules[u.String()] = mod
		d.reverseModules[mod] = u
		return mod, nil
	}

	return nil, err
}

func (d *Document) runModuleScriptSource(script scriptRecord, source string) (p *sobek.Promise, err error) {
	var mod sobek.ModuleRecord

	// run inline scripts directly
	if script.url == nil {
		mod, err = sobek.ParseModule(script.name, source, d.ResolveImportedModule)
		if err != nil {
			return
		}
	} else {
		mod, err = d.fetchModule(script.url)
		if err != nil {
			return
		}
	}

	if err = mod.Link(); err != nil {
		return nil, err
	}

	return mod.Evaluate(d.Realm.VM), nil
}
