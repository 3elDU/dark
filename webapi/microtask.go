package webapi

import (
	"fmt"

	"github.com/grafana/sobek"
)

func installMicrotask(r *Realm) {
	// there's no public method to queue a microtask, but .then() on a promise
	// does exactly that, so we create a promise, resolve it immediately,
	// and call then to queue our custom microtask
	// prototype pollution is not an issue, because this promise is created
	// before any scripts are executed
	p, resolve, _ := r.VM.NewPromise()
	resolve(nil)
	pv := r.VM.ToValue(p)
	then, _ := sobek.AssertFunction(pv.ToObject(r.VM).Get("then"))

	r.VM.Set("queueMicrotask", func(call sobek.FunctionCall) sobek.Value {
		if len(call.Arguments) == 0 {
			panic(r.VM.NewTypeError("At least 1 argument expected, but 0 given"))
		}

		fn, ok := sobek.AssertFunction(call.Arguments[0])
		if !ok {
			panic(r.VM.NewTypeError("Expected first argument to be a function"))
		}

		then(pv, r.VM.ToValue(func(call sobek.FunctionCall) sobek.Value {
			if _, err := fn(sobek.Undefined()); err != nil {
				fmt.Printf("[queueMicrotask] error in microtask: %v\n", err.Error())
				// TODO: report exception as "error" event on window
			}

			return sobek.Undefined()
		}))

		return sobek.Undefined()
	})
}
