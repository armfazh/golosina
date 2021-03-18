// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build amd64,!gccgo,!appengine,call

package curve25519

import (
	"github.com/armfazh/golosina/curve25519/internal/fp"
	"golang.org/x/sys/cpu"
)

var hasBmi2Adx bool

func init() { hasBmi2Adx = cpu.X86.HasBMI2 && cpu.X86.HasADX }

// clamp converts a key into a valid scalar
func clamp(in *[32]byte) *[32]byte {
	out := *in
	out[0] &= 248
	out[31] &= 127
	out[31] |= 64
	return &out
}

// ScalarBaseMult sets dst to the product in*base where dst and base are the x
// coordinates of group points, base is the standard generator and all values
// are in little-endian form.
func ScalarBaseMult(out, in *[32]byte) {
	// The algorithm implemented is the right-to-left Joye's ladder as described
	// in "How to precompute a ladder" in SAC'2017.
	k := clamp(in)
	x1 := &fp.Elt{1} // x1 = 1
	z1 := &fp.Elt{1} // z1 = 1
	x2 := &fp.Elt{   // x2 = G-S
		0xbd, 0xaa, 0x2f, 0xc8, 0xfe, 0xe1, 0x94, 0x7e,
		0xf8, 0xed, 0xb2, 0x14, 0xae, 0x95, 0xf0, 0xbb,
		0xe2, 0x48, 0x5d, 0x23, 0xb9, 0xa0, 0xc7, 0xad,
		0x34, 0xab, 0x7c, 0xe2, 0xee, 0xcd, 0xae, 0x1e}
	z2 := &fp.Elt{1} // z2 = 1
	tt := (*fp.Elt)(out)
	swap := uint(1)
	for s := 0; s < 255-3; s++ {
		i := (s + 3) / 8
		j := (s + 3) % 8
		bit := uint((k[i] >> uint(j)) & 1)
		fp.Cswap(x1, x2, swap^bit)
		fp.Cswap(z1, z2, swap^bit)
		fp.Add(tt, x1, z1)
		fp.Sub(z1, x1, z1)
		fp.Mul(z1, z1, &tableBasePoint255[s])
		fp.Add(x1, tt, z1)
		fp.Sub(z1, tt, z1)
		fp.Sqr(x1, x1)
		fp.Sqr(z1, z1)
		fp.Mul(x1, x1, z2)
		fp.Mul(z1, z1, x2)
		swap = bit
	}
	for i := 0; i < 3; i++ {
		fp.Add(tt, x1, z1)
		fp.Sub(z1, x1, z1)
		fp.Sqr(x1, tt)
		fp.Sqr(z1, z1)
		fp.Sub(x2, x1, z1)
		fp.MulA24(z2, x2)
		fp.Add(z2, z2, z1)
		fp.Mul(x1, x1, z1)
		fp.Mul(z1, x2, z2)
	}
	fp.Inv(z1, z1)
	fp.Mul((*fp.Elt)(out), x1, z1)
	fp.Modp((*fp.Elt)(out))
}

// ScalarMult sets dst to the product in*base where dst and base are the x
// coordinates of group points and all values are in little-endian form.
func ScalarMult(out, in, base *[32]byte) {
	k := clamp(in)
	tt := (*fp.Elt)(out)
	t0, t1 := &fp.Elt{}, &fp.Elt{}
	x1, x2, z2, x3, z3 := &fp.Elt{}, &fp.Elt{}, &fp.Elt{}, &fp.Elt{}, &fp.Elt{}
	*x1 = *base // x1 = xP
	// [RFC-7748] When receiving such an array, implementations
	// of X25519 (but not X448) MUST mask the most significant
	// bit in the final byte.
	x1[31] &= (1 << (255 % 8)) - 1
	x2[0] = 1 // x2 = 1
	*x3 = *x1 // x3 = xP
	z3[0] = 1 // z3 = 1
	move := uint(0)
	for s := 255 - 1; s >= 0; s-- {
		i := s / 8
		j := s % 8
		bit := uint((k[i] >> uint(j)) & 1)
		fp.Add(tt, x2, z2)
		fp.Sub(z2, x2, z2)
		*x2 = *tt
		fp.Add(tt, x3, z3)
		fp.Sub(z3, x3, z3)
		fp.Mul(t0, x2, z3)
		fp.Mul(t1, tt, z2)
		fp.Cmov(x2, tt, move^bit)
		fp.Cmov(z2, z3, move^bit)
		fp.Add(tt, t0, t1)
		fp.Sub(t1, t0, t1)
		fp.Sqr(x3, tt)
		fp.Sqr(z3, t1)
		fp.Mul(z3, x1, z3)
		fp.Sqr(x2, x2)
		fp.Sqr(z2, z2)
		fp.Sub(tt, x2, z2)
		fp.MulA24(t1, tt)
		fp.Add(t1, t1, z2)
		fp.Mul(x2, x2, z2)
		fp.Mul(z2, tt, t1)
		move = bit
	}
	fp.Inv(z2, z2)
	fp.Mul((*fp.Elt)(out), x2, z2)
	fp.Modp((*fp.Elt)(out))
}

//go:noescape
func ladderStep(work *[5]fp.Elt, move uint)

//go:noescape
func difAdd(work *[4]fp.Elt, mu *fp.Elt, swap uint)

//go:noescape
func double(work *[4]fp.Elt)
