package webapi

import (
	"fmt"
	"strings"

	"github.com/grafana/sobek"
)

func installConsole(r *Realm) {
	console := r.VM.NewObject()

	console.Set("log", func(call sobek.FunctionCall) sobek.Value {
		var b strings.Builder
		b.WriteString("[console.log] ")

		for i, arg := range call.Arguments {
			b.WriteString(arg.String())
			if i < len(call.Arguments)-1 {
				b.WriteString(", ")
			}
		}

		fmt.Println(b.String())
		return sobek.Undefined()
	})

	r.VM.Set("console", console)
}
