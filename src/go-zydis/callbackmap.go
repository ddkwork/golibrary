// Copyright 2019 John Papandriopoulos.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package zydis

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

// CallbackMap manages callback tokens in a threadsafe way.
type CallbackMap struct {
	m sync.Map // token uintptr -> value any
	k uintptr
}

// NewToken issues a new token that can be used for a callback.
func (cm *CallbackMap) NewToken(value any) (token unsafe.Pointer) {
	token = unsafe.Pointer(atomic.AddUintptr(&cm.k, 1))
	cm.m.Store(token, value)
	return
}

// GetToken retrieves the value associated with a given token, removes it
// from the map, and returns it.  Panics on an unknown token.
func (cm *CallbackMap) GetToken(token unsafe.Pointer) (value any) {
	value, ok := cm.m.Load(token)
	if !ok {
		panic("CallbackMap: invalid token")
	}
	cm.m.Delete(token)
	return
}
