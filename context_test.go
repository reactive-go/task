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
	"context"
	"testing"

	"github.com/reactive-go/task"
	"github.com/stretchr/testify/assert"
)

func TestGroupDone(t *testing.T) {

	group := task.NewGroup()
	ctx := task.WithGroup(context.Background(), group)
	select {
	case <-ctx.Done():
		assert.Fail(t, "should not be done")
	default:
	}

	group.Signal()
	select {
	case <-ctx.Done():
	default:
		assert.Fail(t, "should be done")
	}
}

func TestGroupFromContext(t *testing.T) {

	ctx := context.Background()
	val, ok := task.GroupFromContext(ctx)
	assert.Nil(t, val)
	assert.False(t, ok)

	group := task.NewGroup()
	ctx = task.WithGroup(ctx, group)
	val, ok = task.GroupFromContext(ctx)
	assert.Equal(t, group, val)
	assert.True(t, ok)
}

func TestOneShotEventDone(t *testing.T) {

	stopEvent := task.NewOneShotEvent()
	ctx := task.WithOneShotEvent(context.Background(), stopEvent)
	select {
	case <-ctx.Done():
		assert.Fail(t, "should not be done")
	default:
	}

	stopEvent.Signal()
	select {
	case <-ctx.Done():
	default:
		assert.Fail(t, "should be done")
	}
}

func TestOneShotEventFromContext(t *testing.T) {

	ctx := context.Background()
	val, ok := task.OneShotEventFromContext(ctx)
	assert.Nil(t, val)
	assert.False(t, ok)

	stopEvent := task.NewOneShotEvent()
	ctx = task.WithOneShotEvent(ctx, stopEvent)
	val, ok = task.OneShotEventFromContext(ctx)
	assert.Equal(t, stopEvent, val)
	assert.True(t, ok)
}
