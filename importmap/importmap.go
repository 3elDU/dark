// Package importmap parses import maps and resolves module specifiers through them.
// https://html.spec.whatwg.org/multipage/webappapis.html#import-maps
package importmap

import (
	"errors"
	"net/url"
)

type ImportMap struct {
	Imports map[string]*url.URL            // nil = blocked
	Scopes  map[string]map[string]*url.URL // scope URL -> specifier map
}

func Parse(text string, baseURL *url.URL) (*ImportMap, error) {
	return nil, errors.ErrUnsupported
}

func (m *ImportMap) Resolve(specifier string, baseURL *url.URL) (*url.URL, error) {
	return nil, errors.ErrUnsupported
}
