package phraseCounter

import (
	"reflect"
	"testing"
)

func TestCountPhraseFrequency(t *testing.T) {
	scenarios := []struct {
		content  string
		expected map[string]int
	}{
		{"", map[string]int{}},
		{"abc", map[string]int{}},
		{"abc def ghi", map[string]int{"abc def ghi": 1}},
		{"abc def ghi hij", map[string]int{"abc def ghi": 1, "def ghi hij": 1}},
		{"abc def.ghi hij", map[string]int{}},
		{"abc def ghi. abc def ghi.", map[string]int{"abc def ghi": 2}},
	}

	for _, s := range scenarios {
		got := CountThreeWordPhraseFrequency(s.content)
		if !reflect.DeepEqual(got, s.expected) {
			t.Errorf("Did not get expected result for content '%v'. Expected %v, got %v\n",
				s.content, s.expected, got)
		}
	}
}
