package markov

import (
	"iter"
	"math/rand/v2"
)

type MarkovState[T comparable] struct {
	MarkovChain[T]
	Current T
	Rand    rand.Rand
}

// Creates some new state
func NewState[T comparable](chain MarkovChain[T], current T, source rand.Source) MarkovState[T] {
	return MarkovState[T]{
		MarkovChain: chain,
		Current:     current,
		Rand:        *rand.New(source),
	}
}

// Returns the next token according to the MarkovState
func (m *MarkovState[T]) GetNext() (T, bool) {
	l, ok := m.MarkovChain.Chain[m.Current]
	if !ok || l.sum == 0 {
		var zero T
		return zero, false
	}

	m.Current = l.GetRand(&m.Rand)

	return m.Current, true
}

// Produces an max sequence of tokens
func (m *MarkovState[T]) Seq() iter.Seq[T] {
	return func(yield func(T) bool) {
		for {
			next, ok := m.GetNext()

			if !ok {
				return
			}

			if !yield(next) {
				return
			}
		}
	}
}
