package ringbuffer

import (
	"sync/atomic"
)

//const (
//	SequenceMax     int64 = (1 << 63) - 1
//	SequenceDefault int64 = -1 // For iinitiializing seq and commit buffer
//	cachePadding               = 7  // Used to align struct to remain in L1 CPU cache
//)

// Cursor is a generic class for tracking a position in a buffer.
// This is for 386x only.
type Cursor struct {
	sequence int64 // Increment value to next available in the series
	//	padding  [cachePadding]int64 // This is probably useless on a 64 bus since align = 8
}

// CursorNew is a factory function that returns a new instance of Cursor
func CursorNew() *Cursor {
	return &Cursor{
		sequence: SequenceDefault,
	}
}

// Store the integer in the sequence.
func (c *Cursor) Store(seq int64) {
	atomic.StoreInt64(&c.sequence, seq)
}

// Load returns the latest values from the sequence.
func (c *Cursor) Load() int64 {
	return atomic.LoadInt64(&c.sequence)
}

// Read returns the latest values from the sequence.
// This is a bit slower than the load function because it allocates 64 bytes on the stack.
func (c *Cursor) Read(noop int64) int64 {
	return atomic.LoadInt64(&c.sequence)
}
