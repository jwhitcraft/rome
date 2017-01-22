package utils

import "sync/atomic"

// Setup a Counter
type Counter int32
func (c *Counter) Increment() int32 {
	return atomic.AddInt32((*int32)(c), 1)
}
func (c *Counter) Reset() {
	atomic.StoreInt32((*int32)(c), 1)
}
func (c *Counter) Get() int32 {
	return atomic.LoadInt32((*int32)(c))
}
// End Counter Setup
