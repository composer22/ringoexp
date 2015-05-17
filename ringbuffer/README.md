## Ringbuffer

A very simple class to manage a ring buffer use between publishers and consumers.

### How it works

One of more publishers use a Leader object to coordinate and place information into a ringbuffer, while one or more consumers use a Follower object to coordinate and process information from a ringbuffer. The Leader and Follower objects validate they do not overwrite each other in the ring by reading each other's commmited status.

### Instructions

To create a ring-buffer, first pre-allocate a specialized array of work you want to track.

For example:
```
ringSize := 1024
ring []*MyWorkStruct := make([]*MyWorkStruct, ringSize)
for i = 0; i < ringSize; i++ {
	ring[i] = &MyWorkStruct{foo: 0}
}
```
Next, we create the components that manage this ring.  This is simply done using this wire up:
```
mngr := ringbuffer.ManagerNew(ringSize)
```
Inside the Manager is a Leader and Follower object

A Leader is used to synchronize the work that needs to be written to the buffer.
A Follower is used to synchonize the work that needs to be consumed from the buffer.
Each is dependent on the other.  The Follower cannot pass the leader.  The Leader cannot
pass the Follower if the Follower has not completed the previous iteration of work.
In essense, each head is chasing the others tail.

Each thread in the system performing work needs to access either a Leader or Follower.

All sizes should be a power of two and be consistent, whether used in Reserver() or in setting up the ring buffer and manager.  If you are using batches, they should be consistent between the Leader and Follower.  If you consistently Reserve 2 slots in the Leader, the Follower should also Reserve 2 slots in its processing.

The reason for this is performance.  Variable batch sizes would need to validate the full range of batch slots cells in the dependent Sequencer. Having boundaries aligned, only the first slot needs to be validated as being available.

Valid:
```
mngr = ringbuffer.ManagerNew(1024)
...
upper = mngr.Leader.Reserve(16)
```

Invalid:
```
mngr = ringbuffer.ManagerNew(500)
...
upper = mngr.Leader.Reserve(12)
...
upper = mngr.Leader.Reserve(14)
```

A publisher thread would be coded something like this to get work into the queue:
```
mask := mngr.Leader.Mask()
batchSize := 8  // Reserve 8 slots of work in the ring
upper := mngr.Leader.Reserve(batchSize)
lower := upper - batchSize + 1
for i := lower; i <= upper; i++ {
	ring[i&mask].foo = 99 // Store some data into the work slot structure
}
mngr.Leader.Commit(lower, upper) // Mark as committed = done.
```
A consumer thread would be coded something like this to work on the queue:
```
mask := mngr.Follower.Mask()
batchSize := 8  // Get 8 slots of work from the ring (same size as Leader)
upper := mngr.Follower.Reserve(batchSize)
lower := upper - batchSize + 1
for i := lower; i <= upper; i++ {
	processWork(ring[i&mask])  // Process something from the ring,
}
mngr.Follower.Commit(lower, upper) // Mark as committed = done.
```
