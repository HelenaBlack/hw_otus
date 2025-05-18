package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type WordCount struct {
	Word  string
	Count int
}

func Top10(text string) []string {
	items := strings.Fields(text)

	wordsMap := make(map[string]int)
	for _, item := range items {
		wordsMap[item]++
	}
	wordCounts := make([]WordCount, 0, len(wordsMap))
	for k, v := range wordsMap {
		wordCounts = append(wordCounts, WordCount{Word: k, Count: v})
	}

	sort.Slice(wordCounts, func(i, j int) bool {
		if wordCounts[i].Count == wordCounts[j].Count {
			return wordCounts[i].Word < wordCounts[j].Word
		}
		return wordCounts[i].Count > wordCounts[j].Count
	})

	if len(wordCounts) > 10 {
		wordCounts = wordCounts[:10]
	}

	wordTop := make([]string, len(wordCounts))
	for i, wc := range wordCounts {
		wordTop[i] = wc.Word
	}

	return wordTop
}
