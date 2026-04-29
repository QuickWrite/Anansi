package markov

import (
	"iter"
	"unicode"
)

// A very basic tokenization that splits the provided string into different string tokens
// by splitting where the whitespace is located. As such a string like
//
//	hello world
//
// would be split into two different tokens in the sequence: "hello" and "world".
//
// The tokenizer also does not differentiate based on the amount of separating whitespace.
// This means that this string is the same to the string before:
//
//	hello           world
//
// Every character that is a space character in terms of unicode is treated as a whitespace character.
func Tokenize(s string) iter.Seq[string] {
	return func(yield func(string) bool) {
		i := skipWhitespace(s, 0)

		for i < len(s) {
			l := i
			for l < len(s) && !unicode.IsSpace(rune(s[l])) {
				l++
			}

			if !yield(s[i:l]) {
				return
			}

			i = skipWhitespace(s, l)
		}
	}
}

func skipWhitespace(s string, pos int) int {
	for pos < len(s) && unicode.IsSpace(rune(s[pos])) {
		pos++
	}

	return pos
}
