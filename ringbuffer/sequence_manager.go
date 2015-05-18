package ringbuffer

// SequenceManager is a simple wrapper around the various compoments for managing a ringbuffer.
type SequenceManager struct {
	Leader   *SeqSingle
	Follower *SeqSingle
}

// SequenceManagerNew instatiates and returns a new Manager.
func SequenceManagerNew(size int64) *SequenceManager {
	m := &SequenceManager{
		Leader:   SeqSingleNew(size, nil, true),
		Follower: SeqSingleNew(size, nil, false),
	}

	// Set the dependencies.
	m.Leader.SetDependency(m.Follower)
	m.Follower.SetDependency(m.Leader)
	return m
}

// SequenceManager is a simple wrapper around the various compoments for managing a ringbuffer.
type SequenceManagerMulti struct {
	Leader   *SeqMulti
	Follower *SeqMulti
}

// SequenceManagerNew instatiates and returns a new Manager.
func SequenceManagerMultiNew(size int64) *SequenceManagerMulti {
	m := &SequenceManagerMulti{
		Leader:   SeqMultiNew(size, nil, true),
		Follower: SeqMultiNew(size, nil, false),
	}

	// Set the dependencies.
	m.Leader.SetDependency(m.Follower)
	m.Follower.SetDependency(m.Leader)
	return m
}

// SequenceManager is a simple wrapper around the various compoments for managing a ringbuffer.
type SequenceManagerX struct {
	Leader   *SeqSimple
	Follower *SeqSimple
}

// SequenceManagerNew instatiates and returns a new Manager.
func SequenceManagerXNew(size int64) *SequenceManagerX {
	m := &SequenceManagerX{
		Leader:   SeqSimpleNew(size, nil, true),
		Follower: SeqSimpleNew(size, nil, false),
	}

	// Set the dependencies.
	m.Leader.SetDependency(m.Follower)
	m.Follower.SetDependency(m.Leader)
	return m
}
