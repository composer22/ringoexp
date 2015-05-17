package ringbuffer

import (
	"math"
	"sync/atomic"
	"time"
)

const (
	SequenceMax     int64 = (1 << 63) - 1
	SequenceDefault int64 = -1 // For iinitiializing seq and commit buffer
)

// Sequencer acts as a reservation hub for tracking and storing any kind of work for a Ring Buffer.
// A single Sequencer can be shared across multiple threads that publish or consume work.
// A publisher is considered a leader; a consumer is considered a follower (behind the leader).
type Sequencer struct {
	cursor     *int64     // Seq number and pointer to the next available pub/con slot.
	dependency *Sequencer // Another sequencer to a pub/con that we are depending on finishing work.
	leader     bool       // Is the sequencer a follower (or a leader?
	buffSize   int64      // The length of the ringbuffer and the committed map.
	committed  []int32    // Tracks the commit states of the work being performed.
	mask       int64      // Used for modulo calculations in indexes.
	shift      uint8      // Used for marking commit states in assignments to commited.
}

// Factory function for returning a new instance of a Sequencer.
func SequencerNew(size int64, dep *Sequencer, leader bool) *Sequencer {
	s := &Sequencer{
		cursor:     new(int64),
		leader:     leader,
		dependency: dep,
		committed:  make([]int32, size),
		buffSize:   size,
		mask:       size - 1,
		shift:      uint8(math.Log2(float64(size))),
	}

	// Init the commit map with default values.
	for i := range s.committed {
		s.committed[i] = int32(SequenceDefault)
	}
	return s
}

// Reserve returns the upper most index for a segment of cells requested by "count".
// For example, if the current pointer is 2 and you needed 5 slots in the ring, the func(5)
// would return 7, the highest cell in the allocation.
// NOTE:
// upper&p.mask is same as (ptr % ring.size) modulo index calculator
// upper >> shift records the number of times that the cursor has been around the ring.
//
func (s *Sequencer) Reserve(count int64) int64 {
	// Set the adjustment value for the barrier.  Assume zero if follower.
	var barAdjust int64
	if s.leader {
		barAdjust = s.buffSize
	}

	// Loop and allocate
	for {
		previous := atomic.LoadInt64(s.cursor) // Get the previous pointer.
		upper := previous + count              // Increment it to get the upper bounds of the chunk.
		lower := previous + 1
		barrier := lower - barAdjust // Calculate the dependency barrier

		// Check dependency based on first cell in the series.  If has not been processed
		// in the last rotation, wait
		for s.dependency.committed[lower&s.mask] != int32(barrier>>s.shift) {
			time.Sleep(time.Microsecond)
		}

		// Update the new sequence number
		if atomic.CompareAndSwapInt64(s.cursor, previous, upper) {
			return upper
		}
	}
}

// Commit updates the committed map to track that a segment in the ring buffer
// has been allocated and used. You call this after calling Reserve() and storing or consuming your
// values into or from the RingBuffer[index from Reserve()]
//
func (s *Sequencer) Commit(lower, upper int64) {
	for ; upper >= lower; upper-- {
		s.committed[upper&s.mask] = int32(upper >> s.shift)
	}
}

// SetDependency is a setter for the dependency of this sequencer.
func (s *Sequencer) SetDependency(d *Sequencer) {
	s.dependency = d
}

// Mask is a getter for the index mask.
func (s *Sequencer) Mask() int64 {
	return s.mask
}
