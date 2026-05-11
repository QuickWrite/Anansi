//   This Source Code Form is subject to the terms of the Mozilla Public
//   License, v. 2.0. If a copy of the MPL was not distributed with this
//   file, You can obtain one at https://mozilla.org/MPL/2.0/.

package markov_test

import (
	"iter"
	"math/rand/v2"
	"testing"

	"github.com/QuickWrite/Anansi/markov"
)

func seqOf[T any](items []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range items {
			if !yield(v) {
				return
			}
		}
	}
}

func TestBuildMarkovChain_Basic(t *testing.T) {
	seq := seqOf([]string{"a", "b", "c"})

	m := markov.BuildMarkovChain(seq)

	r := rand.New(rand.NewPCG(1, 2))

	// a should always go to b
	for range 10 {
		tmp := m.Chain["a"]
		got := tmp.GetRand(r)

		if got != "b" {
			t.Fatalf("expected 'b', got %q", got)
		}
	}

	// b should always go to c
	for range 10 {
		tmp := m.Chain["b"]
		got := tmp.GetRand(r)

		if got != "c" {
			t.Fatalf("expected 'c', got %q", got)
		}
	}
}

func TestBuildMarkovChain_Repeated(t *testing.T) {
	seq := seqOf([]string{"a", "b", "a", "b", "a", "b"})

	m := markov.BuildMarkovChain(seq)

	r := rand.New(rand.NewPCG(3, 4))

	// only possible transition is "b"
	for range 20 {
		tmp := m.Chain["a"]
		got := tmp.GetRand(r)

		if got != "b" {
			t.Fatalf("expected only 'b', got %q", got)
		}
	}
}

func TestBuildMarkovChain_Branching(t *testing.T) {
	seq := seqOf([]string{"a", "b", "a", "c"})

	m := markov.BuildMarkovChain(seq)

	r := rand.New(rand.NewPCG(42, 99))

	counts := map[string]int{}

	for range 1000 {
		tmp := m.Chain["a"]
		v := tmp.GetRand(r)
		counts[v]++
	}

	// both transitions must appear
	if counts["b"] == 0 || counts["c"] == 0 {
		t.Fatalf("expected both 'b' and 'c', got %v", counts)
	}
}

func TestBuildMarkovChain_Single(t *testing.T) {
	seq := seqOf([]string{"a"})

	m := markov.BuildMarkovChain(seq)

	if len(m.Chain) != 0 {
		t.Fatalf("expected no transitions, got %v", m)
	}
}

func TestBuildMarkovChain_LongSequence(t *testing.T) {
	seq := seqOf([]string{"a", "b", "c", "a", "b", "d"})

	m := markov.BuildMarkovChain(seq)

	r := rand.New(rand.NewPCG(7, 8))

	// "b" should sometimes lead to "c" and sometimes "d"
	counts := map[string]int{}

	for range 1000 {
		tmp := m.Chain["b"]
		v := tmp.GetRand(r)
		counts[v]++
	}

	if counts["c"] == 0 || counts["d"] == 0 {
		t.Fatalf("expected both transitions from 'b', got %v", counts)
	}
}
