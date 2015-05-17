package ringbuffer

import "testing"

func TestReserveLeader(t *testing.T) {
	m := ManagerNew(1024)

	// Store a batch of jobs by the Producer.
	for i := 0; i < 8; i++ {
		j := m.Leader.Reserve(1)
		m.Leader.Commit(j, j)
	}

	// Pretend we did the work in the Consumer.
	for i := 0; i < 8; i++ {
		m.Follower.committed[i] = 0
	}

	// Now rotate past beginning again w/ more jobs up to the last cell.
	for i := 0; i < 1024; i++ {
		j := m.Leader.Reserve(1)
		m.Leader.Commit(j, j)
	}
}

func TestReserveFollower(t *testing.T) {
	ringSize := int64(1024)
	m := ManagerNew(ringSize)

	// Pretend the Producer committed a ring of jobs.
	for i := int64(0); i < ringSize; i++ {
		j := m.Leader.Reserve(1)
		m.Leader.Commit(j, j)
	}

	// Now test the follower can read those jobs.
	for i := int64(0); i < ringSize; i++ {
		j := m.Follower.Reserve(1)
		m.Follower.Commit(j, j)
	}

	// Pretend the Producer committed another rotation of jobs.
	for i := int64(0); i < ringSize; i++ {
		j := m.Leader.Reserve(1)
		m.Leader.Commit(j, j)
	}

	// Now test the follower can read another rotation of jobs.
	for i := int64(0); i < ringSize; i++ {
		j := m.Follower.Reserve(1)
		m.Follower.Commit(j, j)
	}

}
