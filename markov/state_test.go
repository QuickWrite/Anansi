package markov_test

import (
	"iter"
	"math/rand/v2"
	"reflect"
	"testing"

	"github.com/QuickWrite/Anansi/markov"
)

func collectN[T any](seq iter.Seq[T], n int) []T {
	var out []T

	for v := range seq {
		out = append(out, v)
		if len(out) == n {
			break
		}
	}

	return out
}

func TestMarkovState_Basic(t *testing.T) {
	seq := seqOf([]string{"a", "b", "c"})
	chain := markov.BuildMarkovChain(seq)

	state := markov.NewState(chain, "a", rand.NewPCG(1, 2))

	next, ok := state.GetNext()
	if !ok || next != "b" {
		t.Fatalf("expected 'b', got %q", next)
	}
}

func TestMarkovState_TerminalStops(t *testing.T) {
	seq := seqOf([]string{"a", "b"})
	chain := markov.BuildMarkovChain(seq)

	state := markov.NewState(chain, "b", rand.NewPCG(1, 2))

	_, ok := state.GetNext()
	if ok {
		t.Fatalf("expected termination, but got value")
	}
}

func TestMarkovState_SeqStops(t *testing.T) {
	seq := seqOf([]string{"a", "b"})
	chain := markov.BuildMarkovChain(seq)

	state := markov.NewState(chain, "a", rand.NewPCG(1, 2))

	var got []string
	for v := range state.Seq() {
		got = append(got, v)
	}

	// should only produce one step: a -> b
	want := []string{"b"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestMarkovState_SequenceProgression(t *testing.T) {
	seq := seqOf([]string{"a", "b", "c", "d"})
	chain := markov.BuildMarkovChain(seq)

	state := markov.NewState(chain, "a", rand.NewPCG(1, 2))

	got := collectN(state.Seq(), 3)

	want := []string{"b", "c", "d"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}
