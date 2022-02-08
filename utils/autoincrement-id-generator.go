package utils

import "sync"

// https://stackoverflow.com/questions/64631848/how-to-create-an-autoincrement-id-field
type AutoInc struct {
	sync.Mutex // ensures AutoInc is goroutine-safe
	id         int
}

func NewAutoInc(startId int) AutoInc {
	return AutoInc{
		id: startId,
	}
}

func (a *AutoInc) ID() (id int) {
	a.Lock()
	defer a.Unlock()

	id = a.id
	a.id++
	return
}
