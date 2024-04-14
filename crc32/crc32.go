// SPDX-License-Identifier: Zlib
//
// Copyright 2024 Andrew Bursavich. All rights reserved.
// Copyright 1995-2024 Jean-loup Gailly and Mark Adler. All rights reserved.
// Use of this source code is governed by the zlib license
// which can be found in the LICENSE file.

// Package crc32 implements the 32-bit cyclic redundancy check, or CRC-32, checksum.
// See https://en.wikipedia.org/wiki/Cyclic_redundancy_check for information.
//
// Polynomials are represented in LSB-first form, also known as reversed representation.
//
// Checksums are layed out in big-endian byte order.
package crc32

import (
	"encoding"
	"hash"
	"hash/crc32"

	"bursavich.dev/crc/internal/shared"
)

// The size of a CRC-32 checksum in bytes.
const Size = 4

var ieeePoly = shared.Val[*Poly]{
	New: func() *Poly {
		return makePoly(crc32.IEEE)
	},
}

// IEEE returns the [Poly] representing the IEEE polynomial,
// which is by far and away the most common CRC-32 polynomial.
// It's used by ethernet (IEEE 802.3), v.42, fddi, gzip, zip, png, ...
func IEEE() *Poly {
	return ieeePoly.Get()
}

var castPoly = shared.Val[*Poly]{
	New: func() *Poly {
		return makePoly(crc32.Castagnoli)
	},
}

// Castagnoli returns the [Poly] representing Castagnoli's polynomial, which is used in iSCSI.
// It has better error detection characteristics than IEEE.
// https://dx.doi.org/10.1109/26.231911
func Castagnoli() *Poly {
	return castPoly.Get()
}

var koopPoly = shared.Val[*Poly]{
	New: func() *Poly {
		return makePoly(crc32.Koopman)
	},
}

// Koopman returns the [Poly] representing Koopman's polynomial.
// It has better error detection characteristics than IEEE.
// https://dx.doi.org/10.1109/DSN.2002.1028931
func Koopman() *Poly {
	return koopPoly.Get()
}

// Hash is a [hash.Hash32] that also implements [encoding.BinaryMarshaler]
// and [encoding.BinaryUnmarshaler] to marshal and unmarshal the internal state
// of the hash. Its Sum methods will lay the value out in big-endian byte order.
type Hash interface {
	hash.Hash32
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

// New creates a new [Hash] computing the CRC-32 checksum using the polynomial
// represented by the [Poly].
func New(p *Poly) Hash {
	return crc32.New(p.stdlib).(Hash)
}

const nBits = Size * 8

// Poly represents a 32-bit polynomial with tables for efficient processing.
type Poly struct {
	poly   uint32
	x2nTbl [nBits]uint32
	stdlib *crc32.Table
}

// MakePoly returns a [Poly] constructed from the specified polynomial
// given in LSB-first form, also known as reversed representation.
// The returned [Poly] may be shared and must not be modified.
func MakePoly(poly uint32) *Poly {
	switch poly {
	case crc32.IEEE:
		return IEEE()
	case crc32.Castagnoli:
		return Castagnoli()
	case crc32.Koopman:
		return Koopman()
	default:
		return makePoly(poly)
	}
}

func makePoly(poly uint32) *Poly {
	p := &Poly{
		poly:   poly,
		stdlib: crc32.MakeTable(poly),
	}
	v := uint32(1) << (nBits - 2)
	p.x2nTbl[0] = v
	for n := 1; n < nBits; n++ {
		v = p.multModP(v, v)
		p.x2nTbl[n] = v
	}
	return p
}

// Polynomial returns the polynomial in LSB-first form, also known as reversed representation.
func (p *Poly) Polynomial() uint32 {
	return p.poly
}

// Checksum returns CRC-32 checksum of the data in big-endian byte order.
func (p *Poly) Checksum(data []byte) uint32 {
	return crc32.Update(0, p.stdlib, data)
}

// Update returns the result of adding the bytes in data to the sum.
func (p *Poly) Update(sum uint32, data []byte) uint32 {
	return crc32.Update(sum, p.stdlib, data)
}

// Combine returns the result of adding n bytes with the next sum to the prev sum.
func (p *Poly) Combine(prev, next uint32, n int64) uint32 {
	if prev == 0 {
		return next
	}
	if n <= 0 {
		return prev
	}
	return p.multModP(prev, p.x2NModP(n, 3)) ^ next
}

// multModP returns (a(x) * b(x)) modulo p(x), where p(x) is the CRC polynomial, reflected.
// For speed, this requires that a is not zero.
func (p *Poly) multModP(a, b uint32) uint32 {
	var v uint32
	for m := uint32(1) << (nBits - 1); m != 0; m >>= 1 {
		if (a & m) != 0 {
			if v ^= b; (a & (m - 1)) == 0 {
				return v
			}
		}
		xor := b&1 != 0
		if b >>= 1; xor {
			b ^= p.poly
		}
	}
	panic("crc32: invalid state")
}

// x2NModP returns x^(n * 2^k) modulo p(x).
func (p *Poly) x2NModP(n int64, k uint32) uint32 {
	v := uint32(1) << (nBits - 1)
	for n != 0 {
		if n&1 != 0 {
			v = p.multModP(p.x2nTbl[k&(nBits-1)], v)
		}
		n >>= 1
		k++
	}
	return v
}
