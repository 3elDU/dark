package browser

import (
	"strings"

	"github.com/3elDU/dark/dom"
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

func classifyScript(scriptNode *dom.Node) scriptType {
	scriptType := strings.ToLower(
		strings.TrimSpace(
			scriptNode.Data.Attributes().Get("type"),
		),
	)

	nomodule := scriptNode.Data.Attributes().Has("nomodule")
	if nomodule {
		// nomodule is for older browsers to load a fallback script.
		// we support modules, so any script with 'nomodule' is ignored.
		return ignoreScript
	}

	switch scriptType {
	// https://mimesniff.spec.whatwg.org/#javascript-mime-type
	case "application/ecmascript", "application/javascript", "application/x-ecmascript",
		"application/x-javascript", "text/ecmascript", "text/javascript", "text/javascript1.0",
		"text/javascript1.1", "text/javascript1.2", "text/javascript1.3", "text/javascript1.4",
		"text/javascript1.5", "text/jscript", "text/livescript", "text/x-ecmascript",
		"text/x-javascript", "":
		return classicScript
	case "module":
		return moduleScript
	case "importmap":
		return importMapScript
	default:
		return ignoreScript
	}
}
