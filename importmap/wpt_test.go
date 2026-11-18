package importmap_test

import (
	"encoding/json"
	"maps"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/3elDU/dark/importmap"
)

// A WPT test object. Children inherit every field their parent sets.
// https://github.com/web-platform-tests/wpt/blob/master/import-maps/data-driven/README.md
type wptTest struct {
	ImportMap               json.RawMessage            `json:"importMap"`
	ImportMapBaseURL        *string                    `json:"importMapBaseURL"`
	BaseURL                 *string                    `json:"baseURL"`
	ExpectedResults         map[string]*string         `json:"expectedResults"`
	ExpectedParsedImportMap json.RawMessage            `json:"expectedParsedImportMap"`
	Tests                   map[string]json.RawMessage `json:"tests"`
}

// ImportMap as WPT serializes it: URLs as strings, null for blocked entries.
type parsedImportMap struct {
	Imports map[string]*string            `json:"imports"`
	Scopes  map[string]map[string]*string `json:"scopes"`
}

func TestWPT(t *testing.T) {
	files, err := filepath.Glob("testdata/wpt/*.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no WPT test files in testdata/wpt")
	}

	for _, file := range files {
		t.Run(strings.TrimSuffix(filepath.Base(file), ".json"), func(t *testing.T) {
			data, err := os.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}

			var root wptTest
			if err := json.Unmarshal(data, &root); err != nil {
				t.Fatal(err)
			}

			runWPTTest(t, root)
		})
	}
}

func runWPTTest(t *testing.T, test wptTest) {
	if test.Tests == nil {
		runLeaf(t, test)
		return
	}

	for _, name := range slices.Sorted(maps.Keys(test.Tests)) {
		var child wptTest
		if err := json.Unmarshal(test.Tests[name], &child); err != nil {
			t.Fatal(err)
		}

		child = inherit(test, child)
		t.Run(name, func(t *testing.T) { runWPTTest(t, child) })
	}
}

func inherit(parent, child wptTest) wptTest {
	if child.ImportMap == nil {
		child.ImportMap = parent.ImportMap
	}
	if child.ImportMapBaseURL == nil {
		child.ImportMapBaseURL = parent.ImportMapBaseURL
	}
	if child.BaseURL == nil {
		child.BaseURL = parent.BaseURL
	}
	if child.ExpectedResults == nil {
		child.ExpectedResults = parent.ExpectedResults
	}
	if child.ExpectedParsedImportMap == nil {
		child.ExpectedParsedImportMap = parent.ExpectedParsedImportMap
	}
	return child
}

func runLeaf(t *testing.T, test wptTest) {
	if test.ImportMap == nil || test.ImportMapBaseURL == nil {
		t.Fatal("test data has no importMap or importMapBaseURL")
	}
	if test.ExpectedParsedImportMap == nil && test.ExpectedResults == nil {
		t.Fatal("test data has no expectations")
	}

	m, err := importmap.Parse(importMapText(test.ImportMap), mustParseURL(t, *test.ImportMapBaseURL))

	if test.ExpectedParsedImportMap != nil {
		checkParsed(t, m, err, test.ExpectedParsedImportMap)
	}
	if test.ExpectedResults != nil {
		checkResolved(t, m, err, test)
	}
}

// A JSON string holds raw import map text, which lets tests feed in invalid JSON.
func importMapText(raw json.RawMessage) string {
	if raw[0] == '"' {
		var text string
		_ = json.Unmarshal(raw, &text)
		return text
	}
	return string(raw)
}

func checkParsed(t *testing.T, m *importmap.ImportMap, parseErr error, expected json.RawMessage) {
	if string(expected) == "null" {
		if parseErr == nil {
			t.Errorf("expected parsing to fail, got %s", toJSON(serialize(m)))
		}
		return
	}

	if parseErr != nil {
		t.Fatalf("parsing failed: %v", parseErr)
	}

	var want parsedImportMap
	if err := json.Unmarshal(expected, &want); err != nil {
		t.Fatal(err)
	}
	want = normalize(want)

	if got := serialize(m); !reflect.DeepEqual(got, want) {
		t.Errorf("parsed import map mismatch\n got: %s\nwant: %s", toJSON(got), toJSON(want))
	}
}

func checkResolved(t *testing.T, m *importmap.ImportMap, parseErr error, test wptTest) {
	if parseErr != nil {
		t.Fatalf("resolution tests need a valid import map, parsing failed: %v", parseErr)
	}
	if test.BaseURL == nil {
		t.Fatal("test data has no baseURL")
	}
	baseURL := mustParseURL(t, *test.BaseURL)

	for _, specifier := range slices.Sorted(maps.Keys(test.ExpectedResults)) {
		want := test.ExpectedResults[specifier]

		t.Run(specifier, func(t *testing.T) {
			got, err := m.Resolve(specifier, baseURL)
			switch {
			case want == nil && err == nil:
				t.Errorf("expected resolution to fail, got %v", got)
			case want != nil && err != nil:
				t.Errorf("resolution failed: %v\nwant: %v", err, *want)
			case want != nil && got.String() != *want:
				t.Errorf("\n got: %v\nwant: %v", got, *want)
			}
		})
	}
}

func serialize(m *importmap.ImportMap) parsedImportMap {
	p := parsedImportMap{Imports: specifierMap(m.Imports)}
	for scope, sm := range m.Scopes {
		if p.Scopes == nil {
			p.Scopes = map[string]map[string]*string{}
		}
		p.Scopes[scope] = specifierMap(sm)
	}
	return normalize(p)
}

func specifierMap(sm map[string]*url.URL) map[string]*string {
	out := make(map[string]*string, len(sm))
	for key, u := range sm {
		if u == nil {
			out[key] = nil
			continue
		}
		s := u.String()
		out[key] = &s
	}
	return out
}

// Treats missing and empty maps as equal.
func normalize(p parsedImportMap) parsedImportMap {
	if p.Imports == nil {
		p.Imports = map[string]*string{}
	}
	if p.Scopes == nil {
		p.Scopes = map[string]map[string]*string{}
	}
	return p
}

func mustParseURL(t *testing.T, s string) *url.URL {
	t.Helper()
	u, err := url.Parse(s)
	if err != nil {
		t.Fatalf("invalid URL in test data %q: %v", s, err)
	}
	return u
}

func toJSON(v any) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}
