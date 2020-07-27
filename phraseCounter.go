package main

import (
	"github.com/dlclark/regexp2"
)

func CountPhraseFrequency(content string) map[string]int {
	phraseFrequency := map[string]int{}

	pattern := `\b(?<first>\w+)\s+(?=((?<second>\w+)(\s+)(?<third>\w+)))`
	re := regexp2.MustCompile(pattern, 0)
	if initialMatch, _ := re.FindStringMatch(content); initialMatch != nil {
		for nextMatch := initialMatch; nextMatch != nil ; nextMatch, _ = re.FindNextMatch(nextMatch) {
			phrase := nextMatch.GroupByName("first").String() + " " +
				nextMatch.GroupByName("second").String() + " " +
				nextMatch.GroupByName("third").String()
			phraseFrequency[phrase] += 1
		}
	}

	return phraseFrequency
}
