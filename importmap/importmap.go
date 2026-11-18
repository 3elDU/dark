// Package importmap parses import maps and resolves module specifiers through them.
// https://html.spec.whatwg.org/multipage/webappapis.html#import-maps
package importmap

import (
	"encoding/json/v2"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type ImportMap struct {
	Imports map[string]*url.URL            // nil = blocked
	Scopes  map[string]map[string]*url.URL // scope URL -> specifier map
}

func parseImportsObject(obj map[string]any, base *url.URL) map[string]*url.URL {
	out := map[string]*url.URL{}

	for k, v := range obj {
		vStr, ok := v.(string)
		if !ok {
			out[k] = nil
			continue
		}

		// if key maps a path prefix, value must end in / too
		// also ignore empty strings
		if k == "" || strings.HasSuffix(k, "/") && !strings.HasSuffix(vStr, "/") {
			out[k] = nil
			continue
		}

		vUrl, err := url.Parse(vStr)
		if err != nil {
			out[k] = nil
			continue
		}

		vUrl = base.ResolveReference(vUrl)
		out[k] = vUrl
	}

	return out
}

func Parse(text string, baseURL *url.URL) (*ImportMap, error) {
	var raw any
	if err := json.Unmarshal([]byte(text), &raw); err != nil {
		fmt.Printf("import map parse error: %v\n", err)
		return nil, err
	}

	rawObj, ok := raw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("importmap is not an object")
	}

	for key := range rawObj {
		if key != "imports" && key != "scopes" && key != "integrity" {
			fmt.Printf("skipping extra '%v' key in importmap\n", key)
		}
	}

	importsObj, ok := rawObj["imports"].(map[string]any)
	// scopesObj, ok := rawObj["scopes"].(map[string]any)

	i := &ImportMap{
		Imports: parseImportsObject(importsObj, baseURL),
	}

	fmt.Println(i.Imports)

	return i, nil
}

func (m *ImportMap) Resolve(specifier string, baseURL *url.URL) (*url.URL, error) {
	return nil, errors.ErrUnsupported
}
