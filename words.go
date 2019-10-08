package words

import (
	"strings"
	"math"
	"sort"
)

type Wordbag struct {
	words map[string]int // per word count
	count int // total words, for performance
}

func NewWordbag() *Wordbag {
	var wb *Wordbag = new(Wordbag)
	wb.words = make(map[string]int)
	return wb
}

// Returns words map, for easy range operations.
func (wb *Wordbag) GetWords() map[string]int {
	return wb.words
}

func (wb *Wordbag) Once(word string) {
	if _, ok := wb.words[word]; !ok {
		wb.words[word] = 1
		wb.count += 1
	}
}

func (wb *Wordbag) None(word string) {
	if _, ok := wb.words[word]; ok {
		wb.count -= wb.words[word]
		delete(wb.words, word)
	}
}

func (wb *Wordbag) Add(word string, count int) {
	if _, ok := wb.words[word]; ok {
		wb.words[word] += count
	} else {
		wb.words[word] = count
	}
	wb.count += count
}

func (wb *Wordbag) Sub(word string, count int) {
	if _, ok := wb.words[word]; ok {
		wb.words[word] -= count
		if wb.words[word] < 1 {
			delete(wb.words, word)
		}
		wb.count -= count
		if wb.count < 0 { // hmm really needed?
			wb.count = 0
		}
	}
}

func (wb *Wordbag) Textract (text string) {
	words := strings.Split(text, " ")
	for _, word := range words {
		wb.Add(word, 1) // filter & normalise words with supplied function ?
	}
}

func (wb *Wordbag) OnceTextract (text string) {
	words := strings.Split(text, " ")
	for _, word := range words {
		wb.Once(word) // filter & normalise words with supplied function ?
	}
}


// just a convenience function...
func (wb *Wordbag) OccurencesTextract (text string) {

	var wordb *Wordbag = NewWordbag()

	wordb.OnceTextract(text)

	wb.OccurencesAdd(wordb)
}

func (wb *Wordbag) Merge(wordb *Wordbag) {
	for w, c := range wordb.words {
		wb.Add(w, c)
	}
}

func (wb *Wordbag) OnceMerge(wordb *Wordbag) {
	for w, _ := range wordb.words {
		wb.Once(w)
	}
}

func  (wb *Wordbag) OccurencesAdd(wordb *Wordbag) {
	for w, _ := range wordb.words {
		wb.Add(w, 1)
	}
}

func (wb *Wordbag) SubMerge(wordb *Wordbag) {
	for w, c := range wordb.words {
		wb.Sub(w, c)
	}
}

func (wb *Wordbag) Clear() {
	for w, _ := range wb.words {
		delete(wb.words, w)
	}
	wb.count = 0
}

func (wb *Wordbag) TotalWords() int {
	return len(wb.words)
}


func (wb *Wordbag) WordCount(w string) int {
	if c, ok := wb.words[w]; ok {
		return c
	}
	return 0
}

// document term count
func (wb *Wordbag) TotalCount() int {
	return wb.count
}


// term frequency
func (wb *Wordbag) TF(w string) float64 {
	if c, ok := wb.words[w]; ok {
		return float64(c)/float64(wb.count)
	}
	return 0.0
}

// inverse document frequency (when wordbag is used to count occurences in corpus)
func (wb *Wordbag) IDF(w string) float64 {
	if c, ok := wb.words[w]; ok {
		return math.Log10(float64(wb.count)/float64(c))
	}
	return 0.0
}

// corpus should be a count of all terms from all documents...
func (wb *Wordbag) Chi2(corpus *Wordbag) float64 {
	var chi float64 = 0.0

	for w, _ := range wb.words {
		E := corpus.TF(w)
		O := wb.TF(w)

		chi += math.Pow(O-E, 2)/E
	}

	return chi
}

func (wb *Wordbag) Top(n int) []string {
	var results []string

	results = make([]string, 0, len(wb.words))

	for k,_ := range wb.words {
		results = append(results, k)
	}

	sort.Slice(results, func(i,j int) bool { return wb.words[results[i]] > wb.words[results[j]] } )

	if c := cap(results); n > c {
		n = c
	}

	if n == 0 {
		// all sorted descending
		return results
	} else {
		return results[0:n]
	}
}

func (wb *Wordbag) Last(n int) []string {
	var results []string

	results = make([]string, 0, len(wb.words))

	for k,_ := range wb.words {
		results = append(results, k)
	}

	sort.Slice(results, func(i,j int) bool { return wb.words[results[i]] < wb.words[results[j]] } )

	if c := cap(results); n > c {
		n = c
	}

	if n == 0 {
		// all sorted ascending
		return results
	} else {
		return results[0:n]
	}
}
