//   This Source Code Form is subject to the terms of the Mozilla Public
//   License, v. 2.0. If a copy of the MPL was not distributed with this
//   file, You can obtain one at https://mozilla.org/MPL/2.0/.

package web

import (
	"embed"
	"hash/fnv"
	"html/template"
	"log"
	"math/rand/v2"
	"net/http"
	"strconv"
	"strings"

	"github.com/QuickWrite/Anansi/markov"
)

// Global variables that store the settings.
var chain markov.MarkovChain[string]
var limit uint
var linkRandomness int
var title string

var tmpl *template.Template
var linkTmpl *template.Template

//go:embed page.gohtml
var page embed.FS

// PageData struct for the go html template of the site
type PageData struct {
	Title   string
	Slug    string
	Content template.HTML
}

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

	var contents strings.Builder

	contents.WriteString(start)

	var i uint = 0
	for p := range s.Seq() {
		i++

		if i > limit {
			break
		}

		contents.WriteString(" ")
		if rand.IntN(linkRandomness) == 1 {
			linkTmpl.Execute(&contents, p)
		} else {
			contents.WriteString(p)
		}
	}

	tmpl.Execute(w, PageData{Title: title, Content: template.HTML(contents.String()), Slug: data})
}

func RunServer(c markov.MarkovChain[string], l, port uint, lr int, t string) {
	chain = c
	limit = l
	linkRandomness = lr
	title = t
	tmpl = template.Must(template.ParseFS(page, "*.gohtml"))
	linkTmpl = tmpl.Lookup("Link")

	http.HandleFunc("/{data}", viewHandler)

	log.Fatal(http.ListenAndServe(":"+strconv.FormatUint(uint64(port), 10), nil))
}
