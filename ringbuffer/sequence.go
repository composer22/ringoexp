package ringbuffer

const (
	SequenceMax     int64 = (1 << 63) - 1
	SequenceDefault int64 = -1 // For iinitiializing seq and commit buffer
)

// Specialized sequencers should implement these methods.
type Sequencer interface {
	Reserve(size int64) int64
	Commit(lower, upper int64)
	SetDependency(d *Sequencer)
	Mask() int64
}
