// SPDX-License-Identifier: Zlib
//
// Copyright 2024 Andrew Bursavich. All rights reserved.
// Use of this source code is governed by the zlib license
// which can be found in the LICENSE file.

package shared

import "sync"

type Val[T any] struct {
	New func() T

	once sync.Once
	val  T
}

func (v *Val[T]) Get() T {
	v.once.Do(func() { v.val = v.New() })
	return v.val
}
