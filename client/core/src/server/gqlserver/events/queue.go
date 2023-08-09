package events

import (
	"sync"
	"time"
)

func newEventsQueue(f func(data interface{}), idleMs int) *eventsQueue {
	result := &eventsQueue{
		queue:  map[string]interface{}{},
		idleMs: idleMs,
		f:      f,
	}
	go result.Start()
	return result
}

type eventsQueue struct {
	queue   map[string]interface{}
	idleMs  int
	f       func(data interface{})
	stopped bool
	sync.Mutex
}

func (eq *eventsQueue) Start() {
	for !eq.stopped {
		eq.Lock()
		queueLen := len(eq.queue)
		if queueLen == 0 {
			eq.Unlock()
			time.Sleep(time.Millisecond * time.Duration(eq.idleMs))
			continue
		}

		data := make([]interface{}, queueLen, queueLen)
		i := 0
		for _, eventData := range eq.queue {
			data[i] = eventData
			i++
		}
		eq.queue = map[string]interface{}{}
		eq.Unlock()
		go eq.f(data)
	}
}

func (eq *eventsQueue) Stop() {
	eq.stopped = true
}

func (eq *eventsQueue) Add(key string, val interface{}) {
	eq.Lock()
	eq.queue[key] = val
	eq.Unlock()
}
