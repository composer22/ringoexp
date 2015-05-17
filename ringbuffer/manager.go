package ringbuffer

// Manager is a simple wrapper around the various compoments for managing a ringbuffer.
type Manager struct {
	Leader   *Sequencer
	Follower *Sequencer
}

// ManagerNew instatiates and returns a new Manager.
func ManagerNew(size int64) *Manager {
	m := &Manager{
		Leader:   SequencerNew(size, nil, true),
		Follower: SequencerNew(size, nil, false),
	}
	// Set the dependencies.
	m.Leader.SetDependency(m.Follower)
	m.Follower.SetDependency(m.Leader)
	return m
}
