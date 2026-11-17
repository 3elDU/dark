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

// Gives access to the higher-level document structure
type Document interface {
	// Needed for dynamic imports by sobek
	ResolveImportedModule(interface{}, string) (sobek.ModuleRecord, error)
	// Needed for import.meta.url
	GetModuleReferrerUrl(any) *url.URL
	// import.meta.resolve()
	ResolveModuleSpecifier(any, string) (*url.URL, error)
}

type Realm struct {
	Url   *url.URL
	Dom   *dom.Dom
	VM    *sobek.Runtime
	Queue *queue.TaskQueue

	Timers         *timers.TimersWorkerHandle
	timerCallbacks map[timers.TimerId]sobek.Callable

	Host Host
	Doc  Document
}

func NewRealm(host Host, doc Document, url *url.URL, dom *dom.Dom) *Realm {
	r := &Realm{
		Url:            url,
		Dom:            dom,
		VM:             sobek.New(),
		Queue:          queue.NewQueue(),
		Timers:         timers.NewTimersWorker(),
		timerCallbacks: map[timers.TimerId]sobek.Callable{},
		Host:           host,
		Doc:            doc,
	}

	r.VM.SetPromiseRejectionTracker(func(p *sobek.Promise, operation sobek.PromiseRejectionOperation) {
		fmt.Printf("unhandled promise rejection: %v\n", p.Result().String())
	})

	installImport(r)
	installWindow(r)
	installTimers(r)
	installConsole(r)
	installMicrotask(r)
	installNavigator(r)
	installDocument(r)

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
