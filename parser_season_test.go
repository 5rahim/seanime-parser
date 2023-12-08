package seanime_parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSeasonAndEpisode(t *testing.T) {

	tests := []struct {
		input    string
		seasons  *[]string
		episodes *[]string
		debug    bool
	}{
		//{"[Seanime] S01 E02 - An episode.mkv", &[]string{"01"}, &[]string{"02"}, true},

		// Single season
		{"[Seanime] Jujutsu Kaisen 2nd Season - 20 [720p][AV1 10bit][AAC][Multi-Sub] (Weekly).mkv", &[]string{"2"}, &[]string{"20"}, false},
		{"[Seanime] Jujutsu Kaisen Season 01.mkv", &[]string{"01"}, nil, false},
		{"[Seanime] Jujutsu Kaisen S1.mkv", &[]string{"1"}, nil, false},
		{"[Seanime] Jujutsu Kaisen 1st Season.mkv", &[]string{"1"}, nil, false},
		{"[Seanime] Jujutsu Kaisen First Season.mkv", &[]string{"1"}, nil, false},
		{"[Seanime] Jujutsu Kaisen S01v2.mkv", &[]string{"01v2"}, nil, false},

		// Season 1 Episode 2
		{"[Seanime] S01E02 - An episode.mkv", &[]string{"01"}, &[]string{"02"}, false},
		{"[Seanime] S01EP02 - An episode.mkv", &[]string{"01"}, &[]string{"02"}, false},
		{"[Seanime] Jujutsu Kaisen 01x02.mkv", &[]string{"01"}, &[]string{"02"}, false},
		{"[Seanime] Jujutsu Kaisen S01E02.mkv", &[]string{"01"}, &[]string{"02"}, false},
		{"[Seanime] Jujutsu Kaisen S1- 02.mkv", &[]string{"1"}, &[]string{"02"}, false},
		{"[Seanime] Jujutsu Kaisen S1-02.mkv", &[]string{"1"}, &[]string{"02"}, false},
		{"[Seanime] Jujutsu Kaisen S1 - 02.mkv", &[]string{"1"}, &[]string{"02"}, false},
		{"[Seanime] Jujutsu Kaisen Season 01 - 02.mkv", &[]string{"01"}, &[]string{"02"}, false},

		{"[Seanime] Jujutsu Kaisen Season 01 - 12.mkv", &[]string{"01"}, &[]string{"12"}, false},

		// Season 1 to 3
		{"[Seanime] Jujutsu Kaisen Seasons 1 ~ 3.mkv", &[]string{"1", "3"}, nil, false},
		{"[Seanime] Jujutsu Kaisen Seasons 01-03.mkv", &[]string{"01", "03"}, nil, false},
		{"[Seanime] Jujutsu Kaisen Season 01-03.mkv", &[]string{"01", "03"}, nil, false},
		{"[Seanime] Jujutsu Kaisen S01-03.mkv", &[]string{"01", "03"}, nil, false},
		{"[Seanime] Jujutsu Kaisen S1-3.mkv", &[]string{"1", "3"}, nil, false},

		// Multiple seasons
		{"[Seanime] Jujutsu Kaisen S1 + S2 + S3.mkv", &[]string{"1", "2", "3"}, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p := newParser(tt.input)
			p.parse()
			if tt.debug {
				t.Log(p.tokenManager.tokens.sDump())
			}

			if tt.seasons != nil {
				assertSeasons(t, p, *tt.seasons)
			} else {
				assertNoSeasons(t, p)
			}

			if tt.episodes != nil {
				assertEpisodes(t, p, *tt.episodes)
			} else {
				assertNoEpisodes(t, p)
			}
		})
	}

}

func TestEpisodeAlt(t *testing.T) {

	tests := []struct {
		input       string
		episodes    *[]string
		episodeAlts *[]string
		debug       bool
	}{
		{"[Seanime] Jujutsu Kaisen 2nd Season - 01 (14) [720p][AV1 10bit][AAC][Multi-Sub] (Weekly).mkv", &[]string{"01"}, &[]string{"14"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p := newParser(tt.input)
			p.parse()
			if tt.debug {
				t.Log(p.tokenManager.tokens.sDump())
			}

			if tt.episodeAlts != nil {
				assertEpisodeAlt(t, p, *tt.episodeAlts)
			} else {
				found, _ := p.tokenManager.tokens.findWithMetadataKind(metadataEpisodeNumberAlt)
				assert.False(t, found)
			}

			if tt.episodes != nil {
				assertEpisodes(t, p, *tt.episodes)
			} else {
				assertNoEpisodes(t, p)
			}
		})
	}

}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func assertSeasons(t *testing.T, p *parser, expectedSeasons []string) {
	found, tkns := p.tokenManager.tokens.findWithMetadataKind(metadataSeason)
	assert.True(t, found)
	assert.Len(t, tkns, len(expectedSeasons))
	for i, tkn := range tkns {
		assert.Equal(t, expectedSeasons[i], tkn.getValue())
	}
}

func assertNoSeasons(t *testing.T, p *parser) {
	found, _ := p.tokenManager.tokens.findWithMetadataKind(metadataSeason)
	assert.False(t, found)
}

func assertEpisodes(t *testing.T, p *parser, expectedEpisodes []string) {
	found, tkns := p.tokenManager.tokens.findWithMetadataKind(metadataEpisodeNumber)
	assert.True(t, found)
	assert.Len(t, tkns, len(expectedEpisodes))
	for i, tkn := range tkns {
		assert.Equal(t, expectedEpisodes[i], tkn.getValue())
	}
}

func assertEpisodeAlt(t *testing.T, p *parser, expectedEpisodes []string) {
	found, tkns := p.tokenManager.tokens.findWithMetadataKind(metadataEpisodeNumberAlt)
	assert.True(t, found)
	assert.Len(t, tkns, len(expectedEpisodes))
	for i, tkn := range tkns {
		assert.Equal(t, expectedEpisodes[i], tkn.getValue())
	}
}

func assertNoEpisodes(t *testing.T, p *parser) {
	found, _ := p.tokenManager.tokens.findWithMetadataKind(metadataEpisodeNumber)
	assert.False(t, found)
}
