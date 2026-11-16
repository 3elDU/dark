package useragent

import (
	"runtime"
	"strings"
)

var userAgents = map[string]string{
	"firefox": "Mozilla/5.0 (X11; Linux x86_64; rv:160.0) Gecko/20100101 Firefox/160.0",
	"chrome":  "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/156.0.0.0 Safari/537.36",
}

var UaOverride string = ""

func OverrideUserAgent(uaKey string) {
	if uaKey != "firefox" && uaKey != "chrome" {
		panic("unknown user agent key")
	}

	UaOverride = uaKey
}

func GetPlatformString() string {
	switch runtime.GOOS {
	case "windows":
		return "Win32"
	case "android", "linux":
		return "Linux " + runtime.GOARCH
	case "darwin", "ios":
		return "MacIntel"
	default:
		return "Win32"
	}
}

func GetAppVersion() string {
	switch runtime.GOOS {
	case "windows":
		return "5.0 (Windows)"
	case "android", "linux":
		return "5.0 (X11)"
	case "darwin", "ios":
		return "5.0 (Macintosh)"
	default:
		return "5.0 (Windows)"
	}
}

func GetDarkVersion() string {
	// TODO: put an actual version here
	return "dark/1.0.0"
}

func GetUserAgent() string {
	if UaOverride != "" {
		return userAgents[UaOverride]
	}

	// assemble a custom useragent
	var sb strings.Builder

	sb.WriteString("Mozilla/5.0 (")
	sb.WriteString(GetPlatformString())
	sb.WriteString(") ")
	sb.WriteString(GetDarkVersion())

	return sb.String()
}
