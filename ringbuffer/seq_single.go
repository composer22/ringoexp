package ringbuffer

import (
	"math"
	"runtime"
)

// SeqSingle is a hub for a single thread/go routine to track access to a ring buffer.
type SeqSingle struct {
	cursor     *int64     // Seq number and pointer to the next available pub/con slot.
	dependency *SeqSingle // Another sequence committed buffer that we are waiting on finishing work.
	leader     bool       // Is the sequence a follower (or a leader?
	buffSize   int64      // The length of the ringbuffer and the committed map.
	committed  []int32    // Tracks the commit states of the work being performed.
	barrier    int64      // Used to calculate downstream or upstream dependencies.
	mask       int64      // Used for modulo calculations in indexes.
	shift      uint8      // Used for marking commit states in assignments to commited.
}

// Factory function for returning a new instance of a SeqSingle.
func SeqSingleNew(size int64, dep *SeqSingle, leader bool) *SeqSingle {
	s := &SeqSingle{
		cursor:     new(int64),
		leader:     leader,
		dependency: dep,
		committed:  make([]int32, size),
		buffSize:   size,
		mask:       size - 1,
		shift:      uint8(math.Log2(float64(size))),
	}

	// Init the cursor and commit map with default values.
	*s.cursor = SequenceDefault
	if leader {
		s.barrier = size
	}

	// Initialize buffer.
	for i := int64(0); i < size; i++ {
		s.committed[i] = int32(SequenceDefault)
	}
	return s
}

// Reserve returns the upper most index for a batch of cells requested by "size".
func (s *SeqSingle) Reserve(size int64) int64 {
	lower := *s.cursor + 1
	gate := lower - s.barrier

	// Check dependency is OK to proceed, based on first cell of that series.
	for s.dependency.committed[lower&s.mask] != int32(gate>>s.shift) {
		runtime.Gosched()
	}
	*s.cursor += size
	return *s.cursor
}

// Reserve returns the upper most index for a batch of cells requested by "size".
func (s *SeqSingle) Reserve2() int64 {
	*s.cursor += 1
	gate := *s.cursor - s.barrier

	// Check dependency is OK to proceed, based on first cell of that series.
	for s.dependency.committed[*s.cursor&s.mask] != int32(gate>>s.shift) {
		runtime.Gosched()
	}

	return *s.cursor
}

// Commit updates the committed map to track that a segment in the ring buffer
// has been allocated and used.
func (s *SeqSingle) Commit(lower, upper int64) {
	for ; upper >= lower; upper-- {
		s.committed[upper&s.mask] = int32(upper >> s.shift)
	}
}

func (s *SeqSingle) Commit2(index int64) {
	s.committed[index&s.mask] = int32(index >> s.shift)
}

// SetDependency is a setter for the dependency of this sequence.
func (s *SeqSingle) SetDependency(d *SeqSingle) {
	s.dependency = d
}

// Mask is a getter for the index mask.
func (s *SeqSingle) Mask() int64 {
	return s.mask
}
