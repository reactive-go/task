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
	"testing"
	"time"

	"github.com/reactive-go/task"
	"github.com/stretchr/testify/assert"
)

type service struct {
	err error
}

func newService(err error) *service {
	return &service{
		err: err,
	}
}

func (rcv *service) Close() error {
	return nil
}

func (rcv *service) Run(stopEv *task.OneShotEvent) error {
	if rcv.err != nil {
		return rcv.err
	}
	<-stopEv.C
	return nil
}

func TestTask_OneGood(t *testing.T) {

	s1 := newService(nil)
	defer s1.Close()

	group := task.NewGroup()
	group.Go(s1)

	select {
	case <-time.After(100 * time.Millisecond):
	case <-group.C:
		assert.Fail(t, "group is signalled")
	}
	assert.NoError(t, group.SignalAndWait())
}

func TestTask_TwoGood(t *testing.T) {

	s1 := newService(nil)
	defer s1.Close()

	s2 := newService(nil)
	defer s2.Close()

	group := task.NewGroup()
	group.Go(s1)
	group.Go(s2)

	select {
	case <-time.After(100 * time.Millisecond):
	case <-group.C:
		assert.Fail(t, "group is signalled")
	}
	assert.NoError(t, group.SignalAndWait())
}

func TestTask_OneBad(t *testing.T) {

	e1 := errors.New("test")
	s1 := newService(e1)
	defer s1.Close()

	group := task.NewGroup()
	group.Go(s1)

	assert.Equal(t, e1, group.Wait())
}

func TestTask_TwoBad(t *testing.T) {

	e1 := errors.New("test")
	s1 := newService(e1)
	defer s1.Close()

	e2 := errors.New("test")
	s2 := newService(e2)
	defer s2.Close()

	group := task.NewGroup()
	group.Go(s1)
	group.Go(s2)

	assert.Equal(t, e1, group.Wait())
}

func TestTask_BadThenGood(t *testing.T) {

	e1 := errors.New("test")
	s1 := newService(e1)
	defer s1.Close()

	s2 := newService(nil)
	defer s2.Close()

	group := task.NewGroup()
	group.Go(s1)
	group.Go(s2)

	assert.Equal(t, e1, group.Wait())
}

func TestTask_GoodThenBad(t *testing.T) {

	s1 := newService(nil)
	defer s1.Close()

	e2 := errors.New("test")
	s2 := newService(e2)
	defer s2.Close()

	group := task.NewGroup()
	group.Go(s1)
	group.Go(s2)

	assert.Equal(t, e2, group.Wait())
}
