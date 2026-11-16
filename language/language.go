package language

import (
	"github.com/jeandeaual/go-locale"
)

var lang string
var langs []string

func init() {
	var err error

	lang, err = locale.GetLocale()
	if err != nil || lang == "" {
		lang = "en-US"
	}

	langs, err = locale.GetLocales()
	if err != nil || len(langs) == 0 {
		langs = []string{"en-US"}
	}
}

func GetUserLanguage() string {
	return lang
}

func GetUserLanguages() []string {
	return langs
}
