package testdata

import "sync"

var lock sync.RWMutex

func DetectRecursiveReadLock() {
	callDoSomething()
}

func callDoSomething() {
	lock.RLock()
	doSomething(5)
	lock.RUnlock()
}

func doSomething(x uint) uint {
	lock.RLock() // want "recursive read lock mutex detected"
	doSomethingMore()
	defer lock.RUnlock()
	return x
}

func doSomethingMore() {
	lock.RLock() // want "recursive read lock mutex detected"
	defer lock.RUnlock()
	return
}
