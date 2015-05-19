package ringbuffer

const (
	SequenceMax     int64 = (1 << 63) - 1
	SequenceDefault int64 = -1 // For iinitiializing seq and commit buffer
)
