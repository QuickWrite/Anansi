//   This Source Code Form is subject to the terms of the Mozilla Public
//   License, v. 2.0. If a copy of the MPL was not distributed with this
//   file, You can obtain one at https://mozilla.org/MPL/2.0/.

package web

import (
	"bytes"
	"embed"
	"hash/fnv"
	"html/template"
	"log"
	"math/rand/v2"
	"net/http"
	"strconv"

	"github.com/QuickWrite/Anansi/markov"
)

// Global variables that store the settings.
var chain markov.MarkovChain[string]
var limit uint
var linkRandomness int
var title string

// PageData struct(s) for the go html template
type PageData struct {
	Title   string
	Slug    string
	Content template.HTML
}

//go:embed page.gohtml
var page embed.FS

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

	contents := ""

	contents += start
	tpl := template.Must(template.ParseFS(page, "*.gohtml"))

	var i uint = 0
	for p := range s.Seq() {
		i++

		if i > limit {
			break
		}

		contents += " "
		if rand.IntN(linkRandomness) == 1 {
			var buf bytes.Buffer
			tpl.ExecuteTemplate(&buf, "Link", p)
			contents += buf.String()
		} else {
			contents += p
		}
	}

	tpl.Execute(w, PageData{Title: title, Content: template.HTML(contents), Slug: data})
}

func RunServer(c markov.MarkovChain[string], l, port uint, lr int, t string) {
	chain = c
	limit = l
	linkRandomness = lr
	title = t

	http.HandleFunc("/{data}", viewHandler)

	log.Fatal(http.ListenAndServe(":"+strconv.FormatUint(uint64(port), 10), nil))
}
