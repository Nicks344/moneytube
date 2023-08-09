package uibridge

import "sync"

type eventHandler func(data interface{}) error

func newEventsContainer() eventsMap {
	return eventsMap{
		handlers: map[string]eventHandler{},
	}
}

var rpcEvents = newEventsContainer()

type eventsMap struct {
	sync.RWMutex
	handlers map[string]eventHandler
}

func (em *eventsMap) add(name string, handler eventHandler) {
	rpcEvents.Lock()
	defer rpcEvents.Unlock()

	em.handlers[name] = handler
}

func (em *eventsMap) get(name string) (eventHandler, bool) {
	rpcEvents.RLock()
	defer rpcEvents.RUnlock()

	res, ok := em.handlers[name]
	return res, ok
}

func (em *eventsMap) remove(name string) {
	rpcEvents.Lock()
	defer rpcEvents.Unlock()

	delete(em.handlers, name)
}
