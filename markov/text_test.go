//   This Source Code Form is subject to the terms of the Mozilla Public
//   License, v. 2.0. If a copy of the MPL was not distributed with this
//   file, You can obtain one at https://mozilla.org/MPL/2.0/.

package markov_test

import (
	"iter"
	"reflect"
	"testing"

	"github.com/QuickWrite/Anansi/markov"
)

// Collects the iterator to a list
func collect[T any](seq iter.Seq[T]) []T {
	var out []T

	for v := range seq {
		out = append(out, v)
	}

	return out
}

func TestTokenize_Word(t *testing.T) {
	got := collect(markov.Tokenize("test"))
	want := []string{"test"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestTokenize_Basic(t *testing.T) {
	got := collect(markov.Tokenize("hello world"))
	want := []string{"hello", "world"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestTokenize_MultipleSpaces(t *testing.T) {
	got := collect(markov.Tokenize("hello     world"))
	want := []string{"hello", "world"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestTokenize_LeadingTrailingWhitespace(t *testing.T) {
	got := collect(markov.Tokenize("   hello world   "))
	want := []string{"hello", "world"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestTokenize_OnlyWhitespace(t *testing.T) {
	got := collect(markov.Tokenize("     \t\n   "))

	if len(got) != 0 {
		t.Fatalf("expected empty slice, got %v", got)
	}
}

func TestTokenize_EmptyString(t *testing.T) {
	got := collect(markov.Tokenize(""))

	if len(got) != 0 {
		t.Fatalf("expected empty slice, got %v", got)
	}
}

func TestTokenize_MixedWhitespace(t *testing.T) {
	input := "hello\tworld\nfoo bar"

	got := collect(markov.Tokenize(input))
	want := []string{"hello", "world", "foo", "bar"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestTokenize_EarlyStop(t *testing.T) {
	seq := markov.Tokenize("one two three")

	var got []string

	for v := range seq {
		got = append(got, v)
		if len(got) == 2 {
			break
		}
	}

	want := []string{"one", "two"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}
