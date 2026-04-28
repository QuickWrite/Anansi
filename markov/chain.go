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

	lower := 0
	higher := len(w.list)

	var i int

	for {
		i = lower + (higher-lower)/2

		// Is the value lower, then go right
		if w.weights[i] < v {
			lower = i
			continue
		}

		// Is i equal to zero                => we reached the end
		// Is the value below greater than v => we reached the end
		if i == 0 || w.weights[i-1] > v {
			return w.list[i]
		}

		higher = i
	}
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
type MarkovChain[T comparable] map[T]WeightedList[T]

// Creates a Markov Chain based on the provided sequence of tokens.
//
// It builds a weighted graph by measuring the frequency of the given states
func BuildMarkovChain[T comparable](seq iter.Seq[T]) MarkovChain[T] {
	markov := MarkovChain[T]{}

	var prev *T = nil

	for elem := range seq {
		if prev == nil {
			prev = &elem
			continue
		}

		if val, ok := markov[*prev]; ok {
			val.addDirty(elem)
		} else {
			markov[*prev] = WeightedList[T]{
				list:    []T{elem},
				weights: []int{1},
				sum:     1,
			}
		}

		prev = &elem
	}

	for _, val := range markov {
		val.clean()
	}

	return markov
}
