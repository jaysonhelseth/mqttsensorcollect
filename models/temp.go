package models

import "sync"

type Temp struct {
	mutex sync.RWMutex
	temp  float64
}

func (t *Temp) Read() float64 {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.temp
}

func (t *Temp) Write(temp float64) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.temp = temp
}
