package main

import "sync/atomic"

// Setup a Counter
type counter int32
func (c *counter) increment() int32 {
	return atomic.AddInt32((*int32)(c), 1)
}
func (c *counter) reset() {
	atomic.StoreInt32((*int32)(c), 1)
}
func (c *counter) get() int32 {
	return atomic.LoadInt32((*int32)(c))
}
// End Counter Setup
