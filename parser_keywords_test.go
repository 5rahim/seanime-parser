package seanime_parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTokenParts(t *testing.T) {

	tests := []struct {
		name              string
		input             string
		expectedTknValues []string
	}{
		{"1", "S01E01", []string{"S", "01", "E", "01"}},
		{"2", "S01E01v2", []string{"S", "01", "E", "01v2"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := newParser(tt.input)
			tm := p.tokenManager

			tkn, ok := tm.tokens.getAtSafe(0)
			assert.True(t, ok)

			// identifyKeyword
			found := p.identifyKeyword(tkn)
			assert.True(t, found)

			tkn, ok = tm.tokens.getAtSafe(0)
			assert.True(t, ok)

			assert.Equal(t, tt.expectedTknValues[0], tkn.getValue())

			tkn, ok = tm.tokens.getAtSafe(1)
			assert.True(t, ok)
			assert.Equal(t, tt.expectedTknValues[1], tkn.getValue())
			assert.True(t, tkn.isMetadataCategory(metadataSeason))

			tkn, ok = tm.tokens.getAtSafe(2)
			assert.True(t, ok)
			assert.Equal(t, tt.expectedTknValues[2], tkn.getValue())

			tkn, ok = tm.tokens.getAtSafe(3)
			assert.True(t, ok)
			assert.Equal(t, tt.expectedTknValues[3], tkn.getValue())
			assert.True(t, tkn.isMetadataCategory(metadataEpisodeNumber))

			t.Log(tm.tokens.sPrint())
		})
	}

}
func TestKeywordGroups(t *testing.T) {

	tests := []struct {
		name             string
		input            string
		expectedTknValue string
		keywordCat       keywordCategory
	}{
		{"1", "BLU RAY 1080P", "BLU RAY", keywordCatSource},
		{"2", "BLU-RAY 1080P", "BLU-RAY", keywordCatSource},
		{"3", "TV RIP 1080P", "TV RIP", keywordCatSource},
		{"4", "10 bits 1080P", "10 bits", keywordCatVideoTerm},
		{"4", "10-bit 1080P", "10-bit", keywordCatVideoTerm},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := newParser(tt.input)
			tm := p.tokenManager

			tkn, _ := tm.tokens.getAtSafe(0)
			found := p.identifyKeyword(tkn)

			assert.True(t, found)
			assert.Equal(t, tt.expectedTknValue, tkn.getValue())
			assert.True(t, tkn.isKeywordCategory(tt.keywordCat))

			t.Log(tm.tokens.sPrint())
		})
	}

}

func TestStandaloneKeywords(t *testing.T) {

	tests := []struct {
		name               string
		input              string
		expectedKeywordCat keywordCategory
	}{
		{"1", "BLURAY 1080P", keywordCatSource},
		{"2", "60FPS 1080P", keywordCatVideoTerm},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := newParser(tt.input)
			tm := p.tokenManager

			tkn, _ := tm.tokens.getAtSafe(0)
			found := p.identifyKeyword(tkn)

			assert.True(t, found)

			tknKeywordCat, found := tkn.getIdentifiedKeywordCategory()

			assert.Equal(t, tt.expectedKeywordCat, tknKeywordCat)

			t.Log(tm.tokens.sPrint())
		})
	}

}
