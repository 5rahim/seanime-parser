package seanime_parser

import (
	"strconv"
	"strings"
)

type parser struct {
	filename     string
	tokenManager *tokenManager
	metadata     *Metadata
}

func newParser(filename string) *parser {
	return &parser{
		filename:     filename,
		tokenManager: newTokenManager(filename),
		metadata:     &Metadata{},
	}
}

func (p *parser) parse() {

	p.parseKeywords()

	p.parseSeason()

}

func (p *parser) parseKeywords() {

	for _, tkn := range *p.tokenManager.tokens {

		// Identify keyword
		_ = p.identifyKeyword(tkn)

	}

}

// identifyKeyword identifies the keyword category of the given token.
func (p *parser) identifyKeyword(tkn *token) bool {

	if tkn.Kind == tokenKindCrc32 {
		tkn.setIdentifiedKeywordCategory(keywordCatFileChecksum)
		return true
	}

	if tkn.Kind == tokenKindPossibleVideoRes {
		tkn.setIdentifiedKeywordCategory(keywordCatVideoResolution)
		return true
	}

	// Check if token is a known pre-defined keyword prefix (e.g. "Blu" for "Blu-ray")
	keywordParts, found := p.tokenManager.keywordManager.findKeywordPartGroups(tkn.getValue())
	foundParts := false
	if found {
		foundParts = false
		for _, keywordGroup := range keywordParts {
			if retTkns, found := p.tokenManager.tokens.peekValuesAfter(p.tokenManager.tokens.getIndexOf(tkn), keywordGroup.seqParts); found {
				// Update token value
				seqPartsStr := ""
				for _, t := range retTkns {
					seqPartsStr += t.getValue()
				}
				tkn.setValue(mergeValues(tkn.getValue(), []string{seqPartsStr}))
				tkn.setIdentifiedKeywordCategory(keywordGroup.category)
				tkn.setKind(tokenKindWord)
				// Remove subsequent tokens
				for _, retTkn := range retTkns {
					p.tokenManager.tokens.removeAt(p.tokenManager.tokens.getIndexOf(retTkn))
				}
				foundParts = true
				break
			}
		}
	}

	if foundParts {
		return true
	}

	// Check if token is a known pre-defined standalone keyword (e.g. "60FPS")
	if len(tkn.getValue()) > 1 {
		if keyword, found := p.tokenManager.keywordManager.findStandaloneKeywordByValue(tkn.getValue()); found {
			tkn.setIdentifiedKeywordCategory(keyword.category)
			return true
		}
	}

	// Parse S01E01
	if strings.HasPrefix(tkn.getNormalizedValue(), "S") && len(tkn.getValue()) > 3 {
		// Extract season and episode
		if season, sep, episode, ok := extractSeasonAndEpisode(tkn.getValue()); ok {
			seasonPrefixTkn := newToken("S")
			seasonPrefixTkn.setIdentifiedKeywordCategory(keywordCatSeasonPrefix)
			seasonPrefixTkn.setKind(tokenKindCharacter)

			seasonTkn := newToken(season)
			seasonTkn.setMetadataCategory(metadataSeason)
			seasonTkn.setKind(tokenKindNumber)

			sepTkn := newToken(sep)
			sepTkn.setIdentifiedKeywordCategory(keywordCatEpisodePrefix)
			sepTkn.setKind(tokenKindCharacter)

			episodeTkn := newToken(episode)
			episodeTkn.setMetadataCategory(metadataEpisodeNumber)
			if isNumber(episode) {
				episodeTkn.setKind(tokenKindNumber)
			} else {
				episodeTkn.setKind(tokenKindNumberLike)
			}

			p.tokenManager.tokens.overwriteAndInsertManyAt(p.tokenManager.tokens.getIndexOf(tkn), []*token{seasonPrefixTkn, seasonTkn, sepTkn, episodeTkn})
			return true
		}
	}

	// Combined or separated seasons
	if strings.HasPrefix(tkn.getNormalizedValue(), "S") {

		if keywords, found := p.tokenManager.keywordManager.findKeywordsBy(func(kw *keyword) bool {
			return kw.isSeasonPrefix() && // Season prefix
				strings.HasPrefix(tkn.getNormalizedValue(), kw.value) // Token starts with season prefix
		}); found {
			for _, keyword := range keywords {

				// e.g. S01
				if keyword.isCombinedWithNumber() {

					// Check if token is after file metadata
					if p.tokenManager.tokens.isTokenAfterFileMetadata(tkn) {
						continue
					}

					// Check if prefix is followed by a number or number-like (e.g. 01, 01v2)
					remaining := strings.TrimPrefix(tkn.getNormalizedValue(), keyword.value)
					if len(remaining) > 0 && isNumberOrLike(remaining) {

						// e.g. S
						seasonPrefixTkn := newToken(keyword.value)
						seasonPrefixTkn.setIdentifiedKeywordCategory(keywordCatSeasonPrefix)
						seasonPrefixTkn.setKind(tokenKindWord)

						// e.g. 01
						seasonTkn := newToken(remaining)
						seasonTkn.setMetadataCategory(metadataSeason)
						seasonTkn.setKind(tokenKindNumberLike)
						if isNumber(remaining) {
							seasonTkn.setKind(tokenKindNumber)
						}

						p.tokenManager.tokens.overwriteAndInsertManyAt(p.tokenManager.tokens.getIndexOf(tkn), []*token{seasonPrefixTkn, seasonTkn})

						return true
					}

				}

				// e.g. Season 01
				if keyword.isSeparatedWithNumber() {

					// Get next token, by skipping delimiters
					// Check if next token is a number or number-like
					if nextTkn, found, _ := p.tokenManager.tokens.getTokenAfterSD(tkn); found &&
						(nextTkn.isNumberOrLikeKind()) {

						nextTkn.setMetadataCategory(metadataSeason)
						return true

					}

				}

				// e.g. 1st Season, first season
				if keyword.isOrdinalSuffix() {

					// Get previous token, by skipping delimiters
					// Check if next token is an ordinal number
					if nextTkn, found, _ := p.tokenManager.tokens.getTokenAfterSD(tkn); found &&
						(nextTkn.isOrdinalNumber()) {

						if num, ok := getNumberFromOrdinal(nextTkn.getValue()); ok {
							nextTkn.setValue(strconv.Itoa(num))
							nextTkn.setMetadataCategory(metadataSeason)
							nextTkn.setKind(tokenKindNumber)
							return true
						}

					}

				}

			}
		}

	}

	return false

}

func (p *parser) parseSeason() {

	for _, tkn := range *p.tokenManager.tokens {

		// Parse 01x01
		if strings.Contains(tkn.getNormalizedValue(), "X") && len(tkn.getValue()) > 3 {
			// Extract season and episode
			if season, sep, episode, ok := extractSeasonAndEpisode(tkn.getValue()); ok {
				seasonTkn := newToken(season)
				seasonTkn.setMetadataCategory(metadataSeason)
				seasonTkn.setKind(tokenKindNumber)

				sepTkn := newToken(sep)
				sepTkn.setIdentifiedKeywordCategory(keywordCatEpisodePrefix)
				sepTkn.setKind(tokenKindCharacter)

				episodeTkn := newToken(episode)
				episodeTkn.setMetadataCategory(metadataEpisodeNumber)
				if isNumber(episode) {
					episodeTkn.setKind(tokenKindNumber)
				} else {
					episodeTkn.setKind(tokenKindNumberLike)
				}

				p.tokenManager.tokens.overwriteAndInsertManyAt(p.tokenManager.tokens.getIndexOf(tkn), []*token{seasonTkn, sepTkn, episodeTkn})
				return
			}
		}

	}

}

// collectMetadata collects the metadata elements from the parsed tokens.
func (p *parser) collectMetadata() {

}
