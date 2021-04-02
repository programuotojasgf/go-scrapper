package phrasecounter

import (
	"github.com/dlclark/regexp2"
)

// CountThreeWordPhraseFrequency returns a map with the occurrence count of 3 word phrases in the given content.
// The key is the phrase and the value is the count of times the phrase was repeated in the given content.
func CountThreeWordPhraseFrequency(content string) map[string]int {
	phraseFrequency := map[string]int{}

	pattern := `\b(?<first>\w+)\s+(?=((?<second>\w+)(\s+)(?<third>\w+)))`
	re := regexp2.MustCompile(pattern, 0)
	if initialMatch, _ := re.FindStringMatch(content); initialMatch != nil {
		for nextMatch := initialMatch; nextMatch != nil; nextMatch, _ = re.FindNextMatch(nextMatch) {
			phrase := nextMatch.GroupByName("first").String() + " " +
				nextMatch.GroupByName("second").String() + " " +
				nextMatch.GroupByName("third").String()
			phraseFrequency[phrase]++
		}
	}

	return phraseFrequency
}
