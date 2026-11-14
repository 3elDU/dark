package browser

type HistoryEntry struct {
	Url string
}

// Represents the history stack within a tab.
type History struct {
	Stack  []HistoryEntry
	active int
}

func NewEmptyHistory() *History {
	return &History{
		Stack: []HistoryEntry{
			{Url: "about:blank"},
		},
		active: 0,
	}
}

func (h *History) Replace(new HistoryEntry) {
	h.Stack[h.active] = new
}

func (h History) Current() HistoryEntry {
	return h.Stack[h.active]
}

func (h *History) Push(new HistoryEntry) {
	h.Stack = append(h.Stack, new)
	h.active = len(h.Stack) - 1
}
