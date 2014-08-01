package bitset_test

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	. "github.com/jrick/bitset"
)

type namedBitSet struct {
	name   string
	bitset BitSet
}

func standardBitsets(numBits uint) []namedBitSet {
	return []namedBitSet{
		{"Words", NewWords(numBits)},
		{"Bytes", NewBytes(numBits)},
		{"Sparse", make(Sparse)},
	}
}

func TestInRange(t *testing.T) {
	tests := []struct {
		name        string
		bitsToSet   []bool
		bitsToUnset []bool // must be same length as bitsToSet
	}{
		{
			bitsToSet:   nil,
			bitsToUnset: nil,
		},
		{
			bitsToSet:   []bool{0: false},
			bitsToUnset: []bool{0: true},
		},
		{
			bitsToSet:   []bool{1: true},
			bitsToUnset: []bool{0: true, 1: true},
		},
		{
			bitsToSet:   []bool{7: true},
			bitsToUnset: []bool{7: true},
		},
		{
			bitsToSet:   []bool{0: true, 7: true},
			bitsToUnset: []bool{1: true, 7: true},
		},
		{
			bitsToSet: []bool{0: true, 1: true, 2: true, 3: true,
				4: true, 5: true, 6: true, 7: false},
			bitsToUnset: []bool{7: true},
		},
		{
			bitsToSet:   []bool{63: true},
			bitsToUnset: []bool{63: true},
		},
		{
			bitsToSet:   []bool{0: true, 8: true, 16: true, 63: true},
			bitsToUnset: []bool{1: true, 9: true, 17: true, 63: true},
		},
		{
			bitsToSet: []bool{56: true, 57: true, 58: true, 59: true,
				60: true, 61: true, 62: true, 63: false},
			bitsToUnset: []bool{63: true},
		},
	}

	for testNum, test := range tests {
	nextBitSet:
		for _, nbs := range standardBitsets(uint(len(test.bitsToSet))) {
			// Set all bits in the bitsToSet field and compare
			// against the expected values.
			for bit, testVal := range test.bitsToSet {
				nbs.bitset.SetBool(uint(bit), testVal)
				got := nbs.bitset.Get(uint(bit))
				if got != testVal {
					t.Errorf("Test %d bitset %s failed: bit %d got %v expected %v",
						testNum, nbs.name, bit, got, testVal)
					continue nextBitSet
				}
			}

			// Unset each bit bit in the bitsToUnset field and
			// check that 1) the value was never set and 2) if
			// unset, the value is now unset.
			for bit, unset := range test.bitsToUnset {
				exp := test.bitsToSet[bit] && !unset
				nbs.bitset.SetBool(uint(bit), exp)
				got := nbs.bitset.Get(uint(bit))
				if got != exp {
					t.Errorf("Test %d bitset %s unset failed: bit %d got %v expected %v",
						testNum, nbs.name, bit, got, exp)
					spew.Dump(nbs.bitset)
					continue nextBitSet
				}
			}
		}
	}
}

type namedGrower struct {
	name   string
	bitset interface {
		BitSet
		Grow(uint)
	}
}

func standardGrowers(numBits uint) []namedGrower {
	words := NewWords(numBits)
	bytes := NewBytes(numBits)
	return []namedGrower{
		{"Words", &words},
		{"Bytes", &bytes},
	}
}

func TestGrowing(t *testing.T) {
	tests := []struct {
		initialBits []bool
		newNumBits  uint
		bitSets     []bool
	}{
		{
			initialBits: nil,
			newNumBits:  1,
			bitSets:     []bool{0: true},
		},
		{
			initialBits: nil,
			newNumBits:  64,
			bitSets:     []bool{0: true, 7: true, 31: true, 63: true},
		},
		{
			initialBits: []bool{0: true, 15: true},
			newNumBits:  64,
			bitSets:     []bool{},
		},
	}

	for testNum, test := range tests {
	nextBitSet:
		for _, nbs := range standardGrowers(uint(len(test.initialBits))) {
			for bit, val := range test.initialBits {
				nbs.bitset.SetBool(uint(bit), val)
			}

			nbs.bitset.Grow(test.newNumBits)

			for bit, val := range test.bitSets {
				nbs.bitset.SetBool(uint(bit), val)
			}

			for bit, exp := range test.bitSets {
				got := nbs.bitset.Get(uint(bit))
				if exp != got {
					t.Errorf("Growing %d bitset %s: bit %d: got %v expected %v",
						testNum, nbs.name, bit, got, exp)
					continue nextBitSet
				}
			}
		}
	}
}

func TestNoSets(t *testing.T) {
	tests := []uint{0, 1, 8, 16, 32, 64, 128, 1024}
	for testNum, test := range tests {
	nextBitSet:
		for _, nbs := range standardBitsets(test) {
			for i := uint(0); i < test; i++ {
				if nbs.bitset.Get(i) {
					t.Errorf("%d: bitset %s: zero value caused set bit %d",
						testNum, nbs.name, i)
					continue nextBitSet
				}
			}
		}
	}
}
