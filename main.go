//   This Source Code Form is subject to the terms of the Mozilla Public
//   License, v. 2.0. If a copy of the MPL was not distributed with this
//   file, You can obtain one at https://mozilla.org/MPL/2.0/.

package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/QuickWrite/Anansi/markov"
	"github.com/QuickWrite/Anansi/web"
)

// Converts the provided arguments into two different buckets
// 1. A list of positional arguments
// 2. A map of flags with their keys and values.
//
// A flag starts with one or two `-` and can contain an argument
// using an = sign. E. g. --flag=argument.
// If there is no provided argument, then the map will contain the
// flag name as an argument and the value as an empty string.
func parseArgs(v []string) ([]string, map[string]string) {
	arguments := []string{}
	flags := map[string]string{}

	for i := range v {
		if v[i][0] == '-' {
			flag := v[i][1:]

			if flag[0] == '-' {
				flag = flag[1:]
			}

			parts := strings.SplitN(flag, "=", 2)

			if len(parts) > 1 {
				flags[parts[0]] = parts[1]
			} else {
				flags[parts[0]] = ""
			}

			continue
		}

		arguments = append(arguments, v[i])
	}

	return arguments, flags
}

func contains[T comparable, V any](m map[T]V, k T) bool {
	_, ok := m[k]

	return ok
}

func printHelp(path string) {
	fmt.Printf(`
Usage:  %s <path> [options]

A mischievous web application to generate an endless stream of senseless text.

  <path>                Path to the text file that will be parsed into a
                        Markov chain.

Options:
  -h, --help            Show this help message and exit.

  -l, --limit=<N>       Maximum number of tokens to generate.
                        If omitted, the default limit is <1000>.

  -p, --port=<PORT>     The port the application should use.
                        If omitted, the default port is <8080>.

  -r=<N>                How often links should be shown on a page.
                        The bigger the number the less links will appear - setting it to 1 or lower disables it.

Examples:
  %s ./books/ulysses.txt                # uses default limit of 1000
  %s ./data/quotes.txt -l=500           # generate up to 500 tokens
  %s ./stories.txt --limit=2000         # generate up to 2000 tokens
  %s ./novel.txt -h                     # display this help screen
  %s ./goethe.txt -r 2                  # will show a lot of links
`, path, path, path, path, path, path)
}

func parseUIntFlag(name string, flags map[string]string) uint {
	l, err := strconv.Atoi(flags[name])

	if err != nil || l <= 0 {
		fmt.Printf("The value for --%s has to be a positive integer >0 and cannot be %s\n", name, flags[name])
		os.Exit(1)
	}

	return uint(l)
}

func main() {
	args, flags := parseArgs(os.Args[1:])

	if contains(flags, "help") || contains(flags, "h") {
		printHelp(os.Args[0])

		os.Exit(0)
	}

	if len(args) != 1 {
		fmt.Printf("The program needs the path to the file.\nTo find out more please write %s --help\n", os.Args[0])
		os.Exit(1)
	}

	file, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}

	var limit uint = 1000
	if contains(flags, "l") {
		limit = parseUIntFlag("l", flags)
	} else if contains(flags, "limit") {
		limit = parseUIntFlag("limit", flags)
	}

	var port uint = 8080
	if contains(flags, "p") {
		port = parseUIntFlag("p", flags)
	} else if contains(flags, "port") {
		port = parseUIntFlag("port", flags)
	}

	var linkRandomness int = 100
	if contains(flags, "r") {
		linkRandomness = int(parseUIntFlag("r", flags))
		if linkRandomness < 1 {
			linkRandomness = 1
		}
	}

	chain := markov.BuildMarkovChain(markov.Tokenize(string(file)))

	log.Print("Running Anansi")

	web.RunServer(chain, limit, port, linkRandomness, "Anansi")
}
