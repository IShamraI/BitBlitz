package locker

import (
	"reflect"
	"sync"
)

const mutexLocked = 1

type Locker struct {
	sync.Mutex
}

func New() *Locker {
	return &Locker{}
}

func (l *Locker) IsLocked() bool {
	state := reflect.ValueOf(l).Elem().FieldByName("state")
	return state.Int()&mutexLocked == mutexLocked
}
