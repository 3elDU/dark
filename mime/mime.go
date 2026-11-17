package mime

// contentType is expected to be trimmed and lowercased
func IsJavascriptContentType(contentType string) bool {
	switch contentType {
	// https://mimesniff.spec.whatwg.org/#javascript-mime-type
	case "application/ecmascript", "application/javascript", "application/x-ecmascript",
		"application/x-javascript", "text/ecmascript", "text/javascript", "text/javascript1.0",
		"text/javascript1.1", "text/javascript1.2", "text/javascript1.3", "text/javascript1.4",
		"text/javascript1.5", "text/jscript", "text/livescript", "text/x-ecmascript",
		"text/x-javascript":
		return true
	default:
		return false
	}
}
