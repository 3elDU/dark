package webapi

import (
	"github.com/3elDU/dark/language"
	"github.com/3elDU/dark/useragent"
	"github.com/grafana/sobek"
)

func installNavigator(r *Realm) {
	nav := r.VM.NewObject()

	nav.Set("appCodeName", "Mozilla")
	nav.Set("appName", "Netscape")
	nav.Set("appVersion", useragent.GetAppVersion())
	nav.Set("platform", useragent.GetPlatformString())
	nav.Set("product", "Gecko")
	nav.Set("productSub", "20100101")
	nav.Set("userAgent", useragent.GetUserAgent())
	nav.Set("vendor", "")
	nav.Set("vendorSub", "")
	nav.Set("oscpu", "Linux x86_64")
	nav.Set("hardwareConcurrency", 2)
	nav.Set("language", language.GetUserLanguage())
	nav.Set("languages", language.GetUserLanguages())
	nav.Set("maxTouchPoints", 0)
	// TODO: run network connectivity checks
	nav.Set("onLine", true)
	nav.Set("pdfViewerEnabled", false)
	nav.Set("plugins", r.VM.NewArray())
	nav.Set("taintEnabled", func(call sobek.FunctionCall) sobek.Value {
		return r.VM.ToValue(false)
	})
	nav.Set("javaEnabled", func(call sobek.FunctionCall) sobek.Value {
		return r.VM.ToValue(false)
	})

	// TODO: missing properties:
	// serviceWorker (Service Worker)
	// storage (Storage API)
	// wakeLock (Wake Lock API)
	// canShare(), share() (Share API)
	// sendBeacon()

	r.VM.Set("navigator", nav)
}
