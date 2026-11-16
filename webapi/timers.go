package webapi

import (
	"math"
	"time"

	"github.com/grafana/sobek"
)

func installTimers(r *Realm) {
	r.VM.Set("setTimeout", func(call sobek.FunctionCall) sobek.Value {
		if len(call.Arguments) == 0 {
			panic(r.VM.NewTypeError("At least 1 argument expected, but 0 passed"))
		}

		fn, ok := sobek.AssertFunction(call.Arguments[0])

		if !ok {
			panic(r.VM.NewTypeError("Expected a function as the first argument, got %v", call.Arguments[0].ExportType().Name()))
		}

		timeout := time.Duration(0)
		if len(call.Arguments) > 1 {
			timeoutMs := call.Arguments[1].ToInteger()

			if timeoutMs < 0 || timeoutMs > int64(math.Pow(2, 32)-1) {
				timeoutMs = 0
			}

			// round to 4ms
			timeoutMs = int64(math.Round(float64(timeoutMs)/4) * 4)

			timeout = time.Duration(timeoutMs * int64(time.Millisecond))
		}

		id := r.Timers.SetTimeout(timeout)
		r.timerCallbacks[id] = fn

		return r.VM.ToValue(r.VM.ToValue(id))
	})

}
