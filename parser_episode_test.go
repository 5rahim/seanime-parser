package seanime_parser

import (
	"testing"
)

func TestOtherEpisodes(t *testing.T) {

	tests := []struct {
		input            string
		expectedTknValue []string
		debug            bool
	}{
		{"[Seanime] Jujutsu Kaisen SP1.mkv", []string{"1"}, false},
		{"[Seanime] Jujutsu Kaisen SP 1.5.mkv", []string{"1.5"}, false},
		{"[Seanime] Jujutsu Kaisen SP1.5.mkv", []string{"1.5"}, true},
		{"[Seanime] Jujutsu Kaisen SP 1.mkv", []string{"1"}, false},
		{"[Seanime] Jujutsu Kaisen OVA 01.mkv", []string{"01"}, false},
		{"[Seanime] Jujutsu Kaisen OVA1.mkv", []string{"1"}, false},
		{"[Seanime] Jujutsu Kaisen NCED1.mkv", []string{"1"}, false},
		{"[Seanime] Jujutsu Kaisen Movie 1.mkv", []string{"1"}, false},
		{"[Seanime] Jujutsu Kaisen Movies 1 ~ 3.mkv", []string{"1", "3"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p := newParser(tt.input)
			p.parse()

			assertOtherEpisodes(t, p, tt.expectedTknValue)

			if tt.debug {
				t.Log(p.tokenManager.tokens.sDump())
			}
		})
	}

}
