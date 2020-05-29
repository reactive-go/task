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
)

// A Group manages a group of goroutines.
//
// A Group combines a OneShotEvent with a WaitGroup,
// so that it can be used to interrupt and initiate the shutdown of one or more long running
// goroutines and then wait for those goroutines to finish.
//
// Groups are typically used as follows:
// - the main goroutine calls Add to set the number of goroutines to wait for;
// - each of the goroutines runs and calls Done when finished;
// - stop is called to initiate the shutdown of the goroutines;
// - wait is called to wait for all goroutines to finished.
//
type Group struct {
	OneShotEvent
	wg  sync.WaitGroup
	mu  sync.Mutex
	err error
}

// NewGroup constructs a new Group instance.
func NewGroup() *Group {
	rcv := &Group{}
	rcv.Init()
	return rcv
}

// Init initialises a Group.
func (rcv *Group) Init() {
	rcv.OneShotEvent.Init()
}

// Add adds delta, which may be negative, to the WaitGroup counter. If the counter becomes zero,
// all goroutines blocked on Wait are released. If the counter goes negative, Add panics.
//
// Note that calls with a positive delta that occur when the counter is zero must happen before a
// Wait. Calls with a negative delta, or calls with a positive delta that start when the counter is
// greater than zero, may happen at any time.
//
// Typically this means the calls to Add should execute before the statement creating the goroutine
// or other event to be waited for.
func (rcv *Group) Add(delta int) {
	rcv.wg.Add(delta)
}

// Done decrements the Group counter by one.
// The first call to Done that specified a non-nil error stops the group;
// its error will be returned by Err and Close.
func (rcv *Group) Done(err error) {
	if err != nil {
		rcv.mu.Lock()
		// Is first non-nil error?
		if rcv.err == nil {
			rcv.err = err
			rcv.mu.Unlock()
			// The first non-nil error stops the group.
			rcv.Signal()
		} else {
			rcv.mu.Unlock()
		}
	}
	rcv.wg.Done()
}

// Err returns the first non-nil error (if any) specified in a call to Done.
// After Err returns a non-nil error, successive calls to Err return the same error.
func (rcv *Group) Err() error {
	rcv.mu.Lock()
	err := rcv.err
	rcv.mu.Unlock()
	return err
}

// SignalAndWait puts the OneShotEvent in a signalled state.
// And then waits until the Group counter is zero.
// Returns the first non-nil error (if any) specified in a call to Done.
func (rcv *Group) SignalAndWait() error {
	rcv.Signal()
	return rcv.Wait()
}

// Wait waits until the Group counter is zero,
// which typically signifies that all function calls from the Go method have returned.
// Returns the first non-nil error (if any) specified in a call to Done.
func (rcv *Group) Wait() error {
	rcv.wg.Wait()
	return rcv.Err()
}

// Go runs the task asynchronously in a new goroutine.
func (rcv *Group) Go(task Task) {
	rcv.Add(1)
	go func() {
		err := task.Run(&rcv.OneShotEvent)
		rcv.Done(err)
	}()
}

// GoFunc calls the function in a new goroutine.
func (rcv *Group) GoFunc(task func(stopEv *OneShotEvent) error) {
	rcv.Go(Func(task))
}
