package limit

import (
	"sync"
)

var Single = newSingle()

// single 防止并发
type single struct {
	locker    sync.Mutex
	callerMap map[string]*caller
}

func newSingle() *single {
	return &single{
		callerMap: make(map[string]*caller),
	}
}

func (s *single) Run(key string, callback func() (interface{}, error)) (interface{}, error) {
	s.locker.Lock()
	c, ok := s.callerMap[key]
	if ok {
		s.locker.Unlock()
		c.waiter.Wait()
		return c.result()
	} else {
		c = newCall()
		s.callerMap[key] = c
		s.locker.Unlock()
	}
	c.run(callback)
	s.locker.Lock()
	delete(s.callerMap, key)
	s.locker.Unlock()
	return c.result()
}

type caller struct {
	value  interface{}
	error  error
	waiter sync.WaitGroup
}

func newCall() *caller {
	c := new(caller)
	c.waiter.Add(1)
	return c
}

func (c *caller) run(fn func() (interface{}, error)) {
	c.value, c.error = fn()
	c.waiter.Done()
}

func (c *caller) result() (interface{}, error) {
	return c.value, c.error
}
