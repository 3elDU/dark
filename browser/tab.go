package browser

import "github.com/3elDU/dark/http"

// A Tab represents a top-level browsing context (a.k.a navigable)
type Tab struct {
	ActiveDocument *Document
	History        *History

	uaKey string
}

func NewBrowsingContext(uaKey string) *Tab {
	t := &Tab{
		History: NewEmptyHistory(),
		uaKey:   uaKey,
	}

	t.ActiveDocument = NewBlankDocument(t)
	return t
}

func (t *Tab) NavigateTo(url string) {
	entry := HistoryEntry{
		Url: url,
	}

	if t.History.Current().Url == "about:blank" {
		t.History.Replace(entry)
	} else {
		t.History.Push(entry)
	}

	page, err := http.FetchWebResource(url, t.uaKey)
	if err != nil {
		panic(err)
	}

	t.ActiveDocument.Close()
	t.ActiveDocument = NewHtmlDocument(t, url, page)
}

func (t *Tab) Close() {
	t.ActiveDocument.Close()
}
