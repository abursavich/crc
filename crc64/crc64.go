// SPDX-License-Identifier: Zlib
//
// Copyright 2024 Andrew Bursavich. All rights reserved.
// Copyright 1995-2024 Jean-loup Gailly and Mark Adler. All rights reserved.
// Use of this source code is governed by the zlib license
// which can be found in the LICENSE file.

// Package crc64 implements the 64-bit cyclic redundancy check, or CRC-64, checksum.
// See https://en.wikipedia.org/wiki/Cyclic_redundancy_check for information.
//
// Polynomials are represented in LSB-first form, also known as reversed representation.
//
// Checksums are layed out in big-endian byte order.
package crc64

import (
	"encoding"
	"hash"
	"hash/crc64"

	"bursavich.dev/crc/internal/lazy"
)

// The size of a CRC-64 checksum in bytes.
const Size = 8

var isoPoly = lazy.Value[*Poly]{
	Init: func() *Poly {
		return makePoly(crc64.ISO)
	},
}

// ISO returns the [Poly] representing the ISO polynomial, defined in ISO 3309 and used in HDLC.
func ISO() *Poly {
	return isoPoly.Get()
}

var ecmaPoly = lazy.Value[*Poly]{
	Init: func() *Poly {
		return makePoly(crc64.ECMA)
	},
}

// ECMA returns the [Poly] representing the ECMA polynomial, defined in ECMA 182.
func ECMA() *Poly {
	return ecmaPoly.Get()
}

// Hash is a [hash.Hash64] that also implements [encoding.BinaryMarshaler]
// and [encoding.BinaryUnmarshaler] to marshal and unmarshal the internal state
// of the hash. Its Sum methods will lay the value out in big-endian byte order.
type Hash interface {
	hash.Hash64
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

// New creates a new [Hash] computing the CRC-64 checksum using the polynomial
// represented by the [Poly].
func New(p *Poly) Hash {
	return crc64.New(p.stdlib).(Hash)
}

const nBits = Size * 8

// Poly represents a 64-bit polynomial with tables for efficient processing.
type Poly struct {
	poly   uint64
	x2nTbl [nBits]uint64
	stdlib *crc64.Table
}

// MakePoly returns a [Poly] constructed from the specified polynomial
// given in LSB-first form, also known as reversed representation.
// The returned [Poly] may be shared and must not be modified.
func MakePoly(poly uint64) *Poly {
	switch poly {
	case crc64.ISO:
		return ISO()
	case crc64.ECMA:
		return ECMA()
	default:
		return makePoly(poly)
	}
}

func makePoly(poly uint64) *Poly {
	p := &Poly{
		poly:   poly,
		stdlib: crc64.MakeTable(poly),
	}
	v := uint64(1) << (nBits - 2)
	p.x2nTbl[0] = v
	for n := 1; n < nBits; n++ {
		v = p.multModP(v, v)
		p.x2nTbl[n] = v
	}
	return p
}

// Polynomial returns the polynomial in LSB-first form, also known as reversed representation.
func (p *Poly) Polynomial() uint64 {
	return p.poly
}

// Checksum returns CRC-64 checksum of the data in big-endian byte order.
func (p *Poly) Checksum(data []byte) uint64 {
	return crc64.Update(0, p.stdlib, data)
}

// Update returns the result of adding the bytes in data to the sum.
func (p *Poly) Update(sum uint64, data []byte) uint64 {
	return crc64.Update(sum, p.stdlib, data)
}

// Combine returns the result of adding n bytes with the next sum to the prev sum.
func (p *Poly) Combine(prev, next uint64, n int64) uint64 {
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
func (p *Poly) multModP(a, b uint64) uint64 {
	var v uint64
	for m := uint64(1) << (nBits - 1); m != 0; m >>= 1 {
		if a&m != 0 {
			if v ^= b; a&(m-1) == 0 {
				return v
			}
		}
		xor := b&1 != 0
		if b >>= 1; xor {
			b ^= p.poly
		}
	}
	panic("crc64: invalid state")
}

// x2NModP returns x^(n * 2^k) modulo p(x).
func (p *Poly) x2NModP(n int64, k uint64) uint64 {
	v := uint64(1) << (nBits - 1)
	for n != 0 {
		if n&1 != 0 {
			v = p.multModP(p.x2nTbl[k&(nBits-1)], v)
		}
		n >>= 1
		k++
	}
	return v
}
