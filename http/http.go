package http

import (
	"fmt"
	"io"
	"net/http"

	"github.com/3elDU/dark/useragent"
	"golang.org/x/net/html/charset"
)

// Used for all requests originating from the browser.
func FetchWebResource(url string) (data string, err error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", useragent.GetUserAgent())

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
