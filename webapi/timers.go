package webapi

import (
	"math"
	"time"

	"github.com/3elDU/dark/timers"
	"github.com/grafana/sobek"
)

func parseIntervalOrTimeout(val sobek.Value) time.Duration {
	duration := time.Duration(0)
	durationMs := val.ToInteger()

	if durationMs < 0 || durationMs > int64(math.Pow(2, 32)-1) {
		durationMs = 0
	}

	// round to 4ms
	durationMs = int64(math.Round(float64(durationMs)/4) * 4)

	duration = time.Duration(durationMs * int64(time.Millisecond))
	return duration
}

func installTimers(r *Realm) {
	r.VM.Set("setTimeout", func(call sobek.FunctionCall) sobek.Value {
		if len(call.Arguments) == 0 {
			panic(r.VM.NewTypeError("At least 1 argument expected, but 0 passed"))
		}

		fn, ok := sobek.AssertFunction(call.Argument(0))
		if !ok {
			panic(r.VM.NewTypeError("Expected a function as the first argument, got %v", call.Arguments[0].ExportType().Name()))
		}

		timeout := parseIntervalOrTimeout(call.Argument(1))
		id := r.Timers.SetTimeout(timeout)
		r.timerCallbacks[id] = fn

		return r.VM.ToValue(id)
	})

	r.VM.Set("setInterval", func(call sobek.FunctionCall) sobek.Value {
		if len(call.Arguments) == 0 {
			panic(r.VM.NewTypeError("At least 1 argument expected, but 0 passed"))
		}

		fn, ok := sobek.AssertFunction(call.Argument(0))
		if !ok {
			panic(r.VM.NewTypeError("Expected a function as the first argument, got %v", call.Arguments[0].ExportType().Name()))
		}

		interval := parseIntervalOrTimeout(call.Argument(1))
		id := r.Timers.SetInterval(interval)
		r.timerCallbacks[id] = fn

		return r.VM.ToValue(id)
	})

	clear := func(id timers.TimerId) {
		r.Timers.Clear(id)
		delete(r.timerCallbacks, id)
	}

	r.VM.Set("clearTimeout", func(call sobek.FunctionCall) sobek.Value {
		id := call.Argument(0).ToInteger()
		clear(timers.TimerId(id))

		return sobek.Undefined()
	})
	r.VM.Set("clearInterval", func(call sobek.FunctionCall) sobek.Value {
		id := call.Argument(0).ToInteger()
		clear(timers.TimerId(id))

		return sobek.Undefined()
	})
}
