package main

import "sync"

// views counts what has been read since the server started. It lives in memory
// on purpose: this is the blog's first piece of shared state, and every
// handler runs in a goroutine of its own, so more than one of them touches it
// at the same moment.
//
// The mutex sits next to what it guards. A counter without it is not merely
// inaccurate: a concurrent write to a map ends the whole process with
// "fatal error: concurrent map writes".
type views struct {
	mu sync.Mutex
	n  map[string]int
}

func newViews() *views {
	return &views{n: make(map[string]int)}
}

// Add counts one reading. Lock lets one goroutine in; defer opens the lock on
// every way out, including a panic — a lock left closed stops not one request
// but all of them.
func (v *views) Add(slug string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.n[slug]++
}

// Count is under the lock too: reading a map while somebody writes to it is
// the same race.
func (v *views) Count(slug string) int {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.n[slug]
}

// Total hands back a copy, so the caller cannot walk the map while a handler
// writes into it.
func (v *views) Total() map[string]int {
	v.mu.Lock()
	defer v.mu.Unlock()

	out := make(map[string]int, len(v.n))
	for slug, n := range v.n {
		out[slug] = n
	}
	return out
}
