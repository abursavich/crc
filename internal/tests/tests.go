// SPDX-License-Identifier: Zlib
//
// Copyright 2024 Andrew Bursavich. All rights reserved.
// Use of this source code is governed by the zlib license
// which can be found in the LICENSE file.

package tests

import (
	"math/rand"
	"testing"
)

type polyCase struct{ a, b []byte }

var polyCases []polyCase

func init() {
	zeroes := make([]byte, 8)
	polyCases = []polyCase{
		{nil, nil},
		{nil, zeroes},
		{zeroes, nil},
		{zeroes, zeroes},
	}
	r := rand.New(rand.NewSource(42))
	for range 128 {
		polyCases = append(polyCases, polyCase{randBuf(r, 256), randBuf(r, 256)})
	}
}

func randBuf(r *rand.Rand, max int) []byte {
	b := make([]byte, r.Intn(max))
	_, _ = r.Read(b)
	return b
}

type PolyFunc func(t *testing.T, a, b []byte)

func FuzzPoly(f *testing.F, fn PolyFunc) {
	for _, c := range polyCases {
		f.Add(c.a, c.b)
	}
	f.Fuzz(fn)
}

func TestPoly(t *testing.T, fn PolyFunc) {
	for _, c := range polyCases {
		c := c
		t.Run("", func(t *testing.T) {
			t.Parallel()
			fn(t, c.a, c.b)
		})
	}
}
