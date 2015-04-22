// Copyright © 2015 Alienero. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"sync"
)

type Map struct {
	m map[interface{}]interface{}
	sync.RWMutex
}

func NewMap() *Map {
	return &Map{
		m:       make(map[interface{}]interface{}),
		RWMutex: sync.RWMutex{},
	}
}

func (m *Map) Add(k, v interface{}) {
	m.Lock()
	m.m[k] = v
	m.Unlock()
}

func (m *Map) Del(k interface{}) {
	m.Lock()
	delete(m.m, k)
	m.Unlock()
}

func (m *Map) Get(k interface{}) interface{} {
	m.RLock()
	defer m.RUnlock()
	return m.m[k]
}
