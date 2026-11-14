package http

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/html/charset"
)

var userAgents = map[string]string{
	"firefox": "Mozilla/5.0 (X11; Linux x86_64; rv:160.0) Gecko/20100101 Firefox/160.0",
	"chrome":  "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/156.0.0.0 Safari/537.36",
	"dark":    "Dark",
}

// Used for all requests originating from the browser.
func FetchWebResource(url string, uaKey string) (data string, err error) {
	if _, ok := userAgents[uaKey]; !ok {
		return "", fmt.Errorf("invalid user agent key: %v", uaKey)
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", userAgents[uaKey])

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	fmt.Printf("response status: %v\n", resp.Status)

	// decode page into a string
	br, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		return
	}

	if body, err := io.ReadAll(br); err == nil {
		return string(body), nil
	} else {
		return "", err
	}
}
