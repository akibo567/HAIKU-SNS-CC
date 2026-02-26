package mora

import (
	"fmt"
	"strings"
	"unicode"
)

// CountMora returns the mora count of the given Japanese text.
// Rules:
//   - Hiragana characters (including small ones like っ, ゃ, ゅ, ょ): 1 mora each
//   - Katakana characters (including small ones like ッ, ャ, ュ, ョ): 1 mora each
//   - Long vowel mark ー: 1 mora
//   - Spaces (full-width and half-width): not counted
//   - Other characters (kanji, latin, numbers): 1 mora each
func CountMora(text string) int {
	count := 0
	for _, r := range text {
		if isSpace(r) {
			continue
		}
		count++
	}
	return count
}

func isSpace(r rune) bool {
	return unicode.IsSpace(r) || r == '　'
}

// ValidateHaiku validates that ku1 has 5 mora, ku2 has 7 mora, ku3 has 5 mora.
func ValidateHaiku(ku1, ku2, ku3 string) error {
	type check struct {
		name string
		text string
		want int
	}

	checks := []check{
		{"上の句", ku1, 5},
		{"中の句", ku2, 7},
		{"下の句", ku3, 5},
	}

	var errs []string
	for _, c := range checks {
		got := CountMora(strings.TrimSpace(c.text))
		if got != c.want {
			errs = append(errs, fmt.Sprintf("%sは%d音でなければなりません（現在: %d音）", c.name, c.want, got))
		}
	}

	if len(errs) > 0 {
		return &ValidationError{Messages: errs}
	}
	return nil
}

type ValidationError struct {
	Messages []string
}

func (e *ValidationError) Error() string {
	return strings.Join(e.Messages, "; ")
}
