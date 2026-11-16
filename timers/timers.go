package timers

import (
	"slices"
	"time"
)

type TimerId uint32

type timer struct {
	id   TimerId
	when time.Time
	// period is set to zero for non-repeating timers
	period time.Duration
}

func (t timer) next() time.Time {
	if t.period == 0 {
		return t.when
	}

	// if this is an interval, compute time for the next tick
	passedTicks := time.Since(t.when) / t.period
	return t.when.Add(t.period * (passedTicks + 1))
}

func closest(timers []timer) (timer, time.Time) {
	var closest time.Time
	var timer timer

	for i, t := range timers {
		deadline := t.next()

		if i == 0 || deadline.Before(closest) {
			closest = deadline
			timer = t
		}
	}

	return timer, closest
}

// push channel used to add new timers, del channel used to delete timers by id
func processTimers(ch chan TimerId, push chan timer, del chan TimerId, stop chan any) {
	timers := []timer{}

	for {
		var fire <-chan time.Time
		var next timer
		if len(timers) > 0 {
			var at time.Time
			next, at = closest(timers)
			fire = time.After(time.Until(at))
		}

		select {
		case t := <-push:
			timers = append(timers, t)
		case id := <-del:
			timers = slices.DeleteFunc(timers, func(t timer) bool { return t.id == id })
		case <-fire:
			select {
			case ch <- next.id:
			case <-stop:
				return
			}

			if next.period == 0 {
				timers = slices.DeleteFunc(timers, func(t timer) bool { return t.id == next.id })
			}
		case <-stop:
			return
		}
	}
}

type TimersWorkerHandle struct {
	// Returns the Timer ID that fired.
	C    chan TimerId
	push chan timer
	del  chan TimerId
	stop chan any

	stopped bool
	next    TimerId
}

// starts a goroutine to process timers
func NewTimersWorker() *TimersWorkerHandle {
	ch := make(chan TimerId)
	push := make(chan timer)
	del := make(chan TimerId)
	stop := make(chan any)

	go processTimers(ch, push, del, stop)

	return &TimersWorkerHandle{
		C:       ch,
		push:    push,
		del:     del,
		stop:    stop,
		stopped: false,
		next:    0,
	}
}

func (h *TimersWorkerHandle) Stop() {
	if h.stopped {
		return
	}

	close(h.stop)
	h.stopped = true
}

func (h *TimersWorkerHandle) SetTimeout(timeout time.Duration) TimerId {
	h.push <- timer{
		id:     h.next,
		when:   time.Now().Add(timeout),
		period: time.Duration(0),
	}

	h.next += 1
	return h.next - 1
}

func (h *TimersWorkerHandle) ClearTimeout(id TimerId) {
	h.del <- id
}

func (h *TimersWorkerHandle) SetInterval(interval time.Duration) TimerId {
	h.push <- timer{
		id:     h.next,
		when:   time.Now(),
		period: interval,
	}

	h.next += 1
	return h.next - 1
}

func (h *TimersWorkerHandle) ClearInterval(id TimerId) {
	h.del <- id
}
