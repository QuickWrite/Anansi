package main

import (
	"hash/fnv"
	"log"
	"math/rand/v2"
	"net/http"

	"github.com/QuickWrite/Anansi/markov"
)

// Global variables that store the settings.
var chain markov.MarkovChain[string]
var limit uint

// NewFromString returns a deterministic RNG seeded from a string.
func NewFromString(seed string) *rand.Rand {
	h := fnv.New64a()
	h.Write([]byte(seed))
	s := h.Sum64()

	return rand.New(rand.NewPCG(s, s^0xdeadbeef))
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	data := r.PathValue("data")

	log.Printf("Received the data: \"%s\"", data)

	rand := NewFromString(data)

	start := chain.GetRandomKey(rand)
	s := markov.NewState(chain, start, rand)
	w.Write([]byte(start))

	var i uint = 0
	for p := range s.Seq() {
		i++

		if i > limit {
			break
		}

		w.Write([]byte(" "))
		w.Write([]byte(p))
	}
}

func runServer(c markov.MarkovChain[string], l uint) {
	chain = c
	limit = l

	http.HandleFunc("/{data}", viewHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
