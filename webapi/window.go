package webapi

import "github.com/grafana/sobek"

func installWindow(r *Realm) {
	g := r.VM.GlobalObject()

	g.DefineAccessorProperty(
		"window",
		r.VM.ToValue(func(call sobek.FunctionCall) sobek.Value {
			return g
		}),
		nil,
		sobek.FLAG_FALSE,
		sobek.FLAG_TRUE,
	)
	r.VM.Set("self", g)
}
