package ringbuffer

// SequenceManager is a simple wrapper around the various compoments for managing a ringbuffer.
type SequenceManager struct {
	Leader   *SeqSimple
	Follower *SeqSimple
}

// SequenceManagerNew instatiates and returns a new Manager.
func SequenceManagerNew(size int64) *SequenceManager {
	m := &SequenceManager{
		Leader:   SeqSimpleNew(size, nil, true),
		Follower: SeqSimpleNew(size, nil, false),
	}

	// Set the dependencies.
	m.Leader.SetDependency(m.Follower)
	m.Follower.SetDependency(m.Leader)
	return m
}
