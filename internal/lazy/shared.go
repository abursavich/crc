// SPDX-License-Identifier: Zlib
//
// Copyright 2024 Andrew Bursavich. All rights reserved.
// Use of this source code is governed by the zlib license
// which can be found in the LICENSE file.

package lazy

import "sync"

type Value[T any] struct {
	Init func() T

	once sync.Once
	val  T
}

func (v *Value[T]) Get() T {
	v.once.Do(func() { v.val = v.Init() })
	return v.val
}
