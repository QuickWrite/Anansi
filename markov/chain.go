//   This Source Code Form is subject to the terms of the Mozilla Public
//   License, v. 2.0. If a copy of the MPL was not distributed with this
//   file, You can obtain one at https://mozilla.org/MPL/2.0/.

package markov

import (
	"iter"
	"math/rand/v2"
)

type WeightedList[T comparable] struct {
	list    []T
	weights []int
	sum     int // The sum of the weights
}

// Returns the next element based on the weights of the list using the provided rand instance.
func (w *WeightedList[T]) GetRand(rand *rand.Rand) T {
	v := rand.IntN(w.sum)

	low, high := 0, len(w.weights)

	for low < high {
		mid := low + (high-low)/2

		if v < w.weights[mid] {
			high = mid
		} else {
			low = mid + 1
		}
	}

	return w.list[low]
}

func (w *WeightedList[T]) addDirty(elem T) {
	for i := 0; i < len(w.list); i++ {
		if w.list[i] != elem {
			continue
		}

		w.weights[i]++

		return
	}

	w.list = append(w.list, elem)
	w.weights = append(w.weights, 1)
}

func (w *WeightedList[T]) clean() {
	sum := 0

	for i := 0; i < len(w.weights); i++ {
		sum += w.weights[i]

		w.weights[i] = sum
	}

	w.sum = sum
}

// Represents a generic Markov Chain to iterate through
type MarkovChain[T comparable] struct {
	Chain map[T]WeightedList[T]
	Keys  []T
}

// Creates a Markov Chain based on the provided sequence of tokens.
//
// It builds a weighted graph by measuring the frequency of the given states
func BuildMarkovChain[T comparable](seq iter.Seq[T]) MarkovChain[T] {
	markov := MarkovChain[T]{
		Chain: map[T]WeightedList[T]{},
		Keys:  []T{},
	}

	var prev *T = nil

	for elem := range seq {
		if prev == nil {
			prev = &elem
			continue
		}

		if val, ok := markov.Chain[*prev]; ok {
			val.addDirty(elem)
			markov.Chain[*prev] = val
		} else {
			markov.Chain[*prev] = WeightedList[T]{
				list:    []T{elem},
				weights: []int{1},
				sum:     1,
			}

			markov.Keys = append(markov.Keys, *prev)
		}

		prev = &elem
	}

	for i, val := range markov.Chain {
		val.clean()
		markov.Chain[i] = val
	}

	return markov
}

func (m *MarkovChain[T]) GetRandomKey(r *rand.Rand) T {
	return m.Keys[r.IntN(len(m.Keys))]
}
