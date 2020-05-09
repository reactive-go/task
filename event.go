// Copyright (C) 2020 Reactive Markets Limited
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package task

import (
	"sync"
	"sync/atomic"
	"time"
)

type (
	// Event is the interface implemented by asynchronous types that may be signalled.
	Event interface {
		Signal()
	}
	// A OneShotEvent is a thread-safe event that can be used to interrupt one or more long
	// running goroutines.
	OneShotEvent struct {
		C         chan struct{}
		signalled int32
	}
)

// NewOneShotEvent constructs a new OneShotEvent.
func NewOneShotEvent() *OneShotEvent {
	rcv := &OneShotEvent{}
	rcv.Init()
	return rcv
}

// Init initialises a OneShotEvent.
func (rcv *OneShotEvent) Init() {
	rcv.C = make(chan struct{})
}

// IsSignalled returns true if the OneShotEvent has been stopped.
func (rcv *OneShotEvent) IsSignalled() bool {
	return atomic.LoadInt32(&rcv.signalled) == 1
}

// Signal puts the OneShotEvent in a signalled state.
// Multiple calls to Signal are idempotent and will not panic.
func (rcv *OneShotEvent) Signal() {
	if atomic.CompareAndSwapInt32(&rcv.signalled, 0, 1) {
		// A closed channel never blocks,
		// so channel closure can be used to signal an "event" to many goroutines.
		close(rcv.C)
	}
}

// Wait blocks the goroutine until the OneShotEvent is signalled.
func (rcv *OneShotEvent) Wait() {
	<-rcv.C
}

// Sleep pauses the current goroutine for the specified duration.
// The function returns false if OneShotEvent was stopped either before or during the the Sleep
// operation.
func (rcv *OneShotEvent) Sleep(d time.Duration) bool {
	t := time.NewTimer(d)
	select {
	case <-t.C:
		return true
	case <-rcv.C:
		// Timer interrupted.
		t.Stop()
		return false
	}
}

// MultiOneShotEvent takes a single stop event and demuxes onto a list of events objects.
type MultiOneShotEvent struct {
	mu        sync.Mutex
	signalled bool
	events    []Event
}

// IsSignalled returns true if the MultiOneShotEvent has been signalled.
func (rcv *MultiOneShotEvent) IsSignalled() (signalled bool) {
	rcv.mu.Lock()
	signalled = rcv.signalled
	rcv.mu.Unlock()
	return
}

// Signal puts the MultiOneShotEvent in a signalled state.
func (rcv *MultiOneShotEvent) Signal() {
	var evs []Event
	rcv.mu.Lock()
	if rcv.signalled {
		// Already signalled.
		rcv.mu.Unlock()
		return
	}
	rcv.signalled = true
	evs = rcv.events
	rcv.events = nil
	rcv.mu.Unlock()

	// Signal after lock has been released.
	for _, ev := range evs {
		ev.Signal()
	}
}

// NewEvent constructs a new OneShotEvent that will be signalled when a shutdown is signalled.
func (rcv *MultiOneShotEvent) NewEvent() *OneShotEvent {
	ev := NewOneShotEvent()
	rcv.Add(ev)
	return ev
}

// NewGroup constructs a new Group that will be signalled when a shutdown is signalled.
func (rcv *MultiOneShotEvent) NewGroup() *Group {
	group := NewGroup()
	rcv.Add(group)
	return group
}

// Add adds a list of events to the MultiOneShotEvent.
// The events will be signalled immediately if the MultiOneShotEvent has already been signalled.
func (rcv *MultiOneShotEvent) Add(evs ...Event) {
	rcv.mu.Lock()
	if !rcv.signalled {
		rcv.events = append(rcv.events, evs...)
		rcv.mu.Unlock()
		return
	}
	// Signal immediately if the MultiOneShotEvent has already been signalled.
	for _, ev := range evs {
		ev.Signal()
	}
}
