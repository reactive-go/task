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

type (
	// Task is the interface implemented by all task types.
	Task interface {
		Run(stopEv *OneShotEvent) error
	}

	// Func is the function type for the Run method.
	Func func(stopEv *OneShotEvent) error
)

// Run implements the Task interface by delegating to a Task function.
func (rcv Func) Run(stopEv *OneShotEvent) error {
	return rcv(stopEv)
}
