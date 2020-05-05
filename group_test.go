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

package task_test

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/reactive-go/task"
	"github.com/stretchr/testify/assert"
)

func TestGroup(t *testing.T) {
	group := task.NewGroup()
	assert.False(t, group.IsSignalled())
	assert.True(t, group.Sleep(10*time.Millisecond))

	errTest := errors.New("test")
	group.Add(1)
	go func() {
		defer group.Done(errTest)
		<-group.C
	}()
	group.Signal()
	assert.True(t, group.IsSignalled())
	group.Signal() // Idempotent.
	assert.True(t, group.IsSignalled())
	assert.False(t, group.Sleep(10*time.Second))

	group.Wait()
	assert.Equal(t, errTest, group.Err())
}

func TestWaitGroup(t *testing.T) {
	var count int32
	stop := make(chan struct{})
	done := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		done.Add(1)
		go func() {
			// Do not increment until stop channel is closed.
			<-stop
			atomic.AddInt32(&count, 1)
			done.Done()
		}()
	}
	assert.Equal(t, int32(0), atomic.LoadInt32(&count))
	// Signal the stop event.
	close(stop)
	done.Wait()
	assert.Equal(t, int32(10), atomic.LoadInt32(&count))
}

func BenchmarkIsStopped(b *testing.B) {
	group := task.NewGroup()
	for i := 0; i < b.N; i++ {
		if group.IsSignalled() {
			return
		}
	}
}

func BenchmarkSelectStopped(b *testing.B) {
	group := task.NewGroup()
	for i := 0; i < b.N; i++ {
		select {
		case <-group.C:
			return
		default:
		}
	}
}
