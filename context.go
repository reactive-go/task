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
	"context"
)

type key int

const (
	groupKey key = iota + 1
	oneShotEvent
)

// WithGroup returns a new Context that holds the Group.
func WithGroup(parent context.Context, val *Group) context.Context {
	return &groupCtx{
		Context: parent,
		group:   val,
	}
}

// GroupFromContext returns the Group stored in the Context.
func GroupFromContext(ctx context.Context) (*Group, bool) {
	val, ok := ctx.Value(groupKey).(*Group)
	return val, ok
}

type groupCtx struct {
	context.Context
	group *Group
}

func (rcv *groupCtx) Done() <-chan struct{} {
	return rcv.group.C
}

// Err returns context.Canceled if the Done channel is closed or nil otherwise.
// After Err returns a non-nil error, successive calls to Err return the same error.
// The Err method is overridden to maintain the context.Canceled contract.
func (rcv *groupCtx) Err() error {
	if rcv.group.IsSignalled() {
		return context.Canceled
	}
	return nil
}

func (rcv *groupCtx) Value(key interface{}) interface{} {
	if key == groupKey {
		return rcv.group
	}
	return rcv.Context.Value(key)
}

// WithOneShotEvent returns a new Context that holds the Event.
func WithOneShotEvent(parent context.Context, val *OneShotEvent) context.Context {
	return &oneShotEventCtx{
		Context: parent,
		event:   val,
	}
}

// OneShotEventFromContext returns the OneShotEvent stored in the Context.
func OneShotEventFromContext(ctx context.Context) (*OneShotEvent, bool) {
	val, ok := ctx.Value(oneShotEvent).(*OneShotEvent)
	return val, ok
}

type oneShotEventCtx struct {
	context.Context
	event *OneShotEvent
}

func (rcv *oneShotEventCtx) Done() <-chan struct{} {
	return rcv.event.C
}

// Err returns context.Canceled if the Done channel is closed or nil otherwise.
// After Err returns a non-nil error, successive calls to Err return the same error.
// The Err method is overridden to maintain the context.Canceled contract.
func (rcv *oneShotEventCtx) Err() error {
	if rcv.event.IsSignalled() {
		return context.Canceled
	}
	return nil
}

func (rcv *oneShotEventCtx) Value(key interface{}) interface{} {
	if key == oneShotEvent {
		return rcv.event
	}
	return rcv.Context.Value(key)
}
