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
	"testing"
	"time"

	"github.com/reactive-go/task"
	"github.com/stretchr/testify/assert"
)

func TestOneShotEvent(t *testing.T) {
	ev := task.NewOneShotEvent()
	assert.False(t, ev.IsSignalled())
	assert.True(t, ev.Sleep(10*time.Millisecond))

	go func() {
		<-ev.C
	}()

	ev.Signal()
	assert.True(t, ev.IsSignalled())
	ev.Signal() // Idempotent.
	assert.True(t, ev.IsSignalled())
	assert.False(t, ev.Sleep(10*time.Second))

	ev.Wait()
}

func TestMultiOneShotEvent(t *testing.T) {
	ev := task.MultiOneShotEvent{}
	s1 := task.NewOneShotEvent()
	s2 := task.NewOneShotEvent()
	ev.Add(s1, s2)
	assert.False(t, ev.IsSignalled())
	assert.False(t, s1.IsSignalled())
	assert.False(t, s2.IsSignalled())

	ev.Signal()
	assert.True(t, ev.IsSignalled())
	assert.True(t, s1.IsSignalled())
	assert.True(t, s2.IsSignalled())

	s3 := task.NewOneShotEvent()
	ev.Add(s3)
	assert.True(t, s3.IsSignalled())
}
