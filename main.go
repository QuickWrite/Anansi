package main

import (
	"fmt"
	"log"
	"os"

	"github.com/QuickWrite/Anansi/markov"
)

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

	runServer(chain, 1000)
}
