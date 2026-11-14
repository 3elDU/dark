package browser

import "github.com/3elDU/dark/http"

// A browsing context (a.k.a navigable) represents a top-level browsing context.
type Tab struct {
	ActiveDocument *Document
	History        *History

	uaKey string
}

func NewBrowsingContext(uaKey string) *Tab {
	return &Tab{
		ActiveDocument: NewBlankDocument(),
		History:        NewEmptyHistory(),
		uaKey:          uaKey,
	}
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

	t.ActiveDocument = NewHtmlDocument(page)
}
