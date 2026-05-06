package main

import (
	"fmt"
	"hash/fnv"
	"log"
	"math/rand/v2"
	"net/http"
	"os"

	"github.com/QuickWrite/Anansi/markov"
)

var chain markov.MarkovChain[string]

const limit = 1000

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

	i := 0
	for p := range s.Seq() {
		i++

		if i > limit {
			break
		}

		w.Write([]byte(" "))
		w.Write([]byte(p))
	}
}

func main() {
	args := os.Args[1:]

	if len(args) != 1 {
		fmt.Printf("The program needs the path to the file.")
		os.Exit(1)
	}

	file, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}

	chain = markov.BuildMarkovChain(markov.Tokenize(string(file)))

	log.Print("Running Anansi")

	http.HandleFunc("/{data}", viewHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
