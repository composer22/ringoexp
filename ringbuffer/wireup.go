package ringbuffer

type Manager struct {
	Leader   *Sequencer
	Follower *Sequencer
}

// ManagerNew instatiates and returns a new Manager
func ManagerNew(size int64) *Manager {
	m := &Manager{
		Leader:   SequencerNew(size, nil, true),
		Follower: SequencerNew(size, nil, false),
	}
	m.Leader.SetDependency(m.Follower)
	m.Follower.SetDependency(m.Leader)
	return m
}
