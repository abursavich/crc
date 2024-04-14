// SPDX-License-Identifier: Zlib
//
// Copyright 2024 Andrew Bursavich. All rights reserved.
// Use of this source code is governed by the zlib license
// which can be found in the LICENSE file.

package crc32

import (
	"hash/crc32"
	"math/bits"
	"testing"

	"bursavich.dev/crc/internal/tests"
)

func FuzzPoly(f *testing.F) {
	tests.FuzzPoly(f, testPoly)
}

func TestPoly(t *testing.T) {
	tests.TestPoly(t, testPoly)
}

var polys = []*Poly{
	MakePoly(crc32.IEEE),
	MakePoly(crc32.Castagnoli),
	MakePoly(crc32.Koopman),
	MakePoly(bits.Reverse32(crc32.Castagnoli)),
}

func testPoly(t *testing.T, a, b []byte) {
	for _, p := range polys {
		aSum := p.Checksum(a)
		bSum := p.Checksum(b)
		want := p.Update(aSum, b)

		h := New(p)
		h.Write(a)
		h.Write(b)
		if got := h.Sum32(); got != want {
			t.Errorf("Poly = 0x%08x; Hash.Sum32() = 0x%08x; want 0x%08x", p.poly, got, want)
		}

		if got := p.Combine(aSum, bSum, int64(len(b))); got != want {
			t.Errorf("Poly = 0x%08x; Combine(0x%08x, 0x%08x, %d) = 0x%08x; want 0x%08x", p.poly, aSum, bSum, len(b), got, want)
		}
	}
}
