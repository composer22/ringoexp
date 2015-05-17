## Ringbuffer Instructions

To create a ring-buffer, first allocate an array of work you want to track.

For example:
```
ringSize := 1024
ring []*MyWorkStruct := make([]*MyWorkStruct, ringSize)
for i = 0; i < ringSize; i++ {
	ring[i] = &MyWorkStruct{foo: 0}
}
```
Next, we create the components that manage this ring.  This is simply done:
```
mngr := ringbuffer.ManagerNew(ringSize)
```
A Leader is used to synchronize the work that needs to be written to the buffer.
A Follower is used to synchonize the work that needs to be consumed from the buffer.
Each is dependent on the other.  The Follower cannot pass the leader.  The Leader cannot
pass the Follower if the Follower has not completed the previous iteration of work.
In essense, each head is chasing the others tail.

Each thread in the system performing work needs to access either a Leader or Follower.

A publisher thread would perform like this to get work into the queue:
```
mask := mngr.Leader.Mask()
size := 6  // Reserve 6 slots of work in the ring
upper := mngr.Leader.Reserve(size)
lower := upper - size + 1
for i := lower; i <= upper; i++ {
	ring[i&mask].foo = 99 // Store some data into the work slot structure
}
mngr.Leader.Commit(lower, upper) // Mark as committed = done.
```
A consumer thread would perform like this to work on the queue:
```
mask := mngr.Follower.Mask()
size := 6  // Get 6 slots of work from the ring
upper := mngr.Follower.Reserve(size)
lower := upper - size + 1
for i := lower; i <= upper; i++ {
	processWork(ring[i&mask])  // Process something from the ring,
}
mngr.Follower.Commit(lower, upper) // Mark as committed = done.
```
