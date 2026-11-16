// Package webapi provides implementation of all the Web APIs for
// JavaScript code running inside the webpage.
package webapi

import (
	"fmt"
	"net/url"

	"github.com/3elDU/dark/dom"
	"github.com/3elDU/dark/queue"
	"github.com/3elDU/dark/timers"
	"github.com/grafana/sobek"
)

// Host gives access to the higher-level structure that holds this realm
// (a tab)
type Host interface {
	NavigateTo(page string)
}

type Realm struct {
	Url   *url.URL
	Dom   *dom.Dom
	VM    *sobek.Runtime
	Queue *queue.TaskQueue

	Timers         *timers.TimersWorkerHandle
	timerCallbacks map[timers.TimerId]sobek.Callable

	Host Host
}

func NewRealm(host Host, url *url.URL, dom *dom.Dom) *Realm {
	r := &Realm{
		Url:            url,
		Dom:            dom,
		VM:             sobek.New(),
		Queue:          queue.NewQueue(),
		Timers:         timers.NewTimersWorker(),
		timerCallbacks: map[timers.TimerId]sobek.Callable{},
		Host:           host,
	}

	installWindow(r)
	installTimers(r)
	installConsole(r)
	installMicrotask(r)

	return r
}

func (r *Realm) Close() {
	r.Timers.Stop()
}

func (r *Realm) MainLoop() {
	for {
		select {
		case id := <-r.Timers.C:
			fn, ok := r.timerCallbacks[id]
			if !ok {
				fmt.Printf("[Realm] callback not found for timer %v\n", id)
				break
			}

			r.Queue.ScheduleFunc(func(vm *sobek.Runtime) error {
				_, err := fn(nil)
				return err
			})
		}

		r.Queue.DrainQueue(r.VM)
	}
}
