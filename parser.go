package seanime_parser

import "strings"

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
	if strings.HasPrefix(tkn.getValue(), "S") && len(tkn.getValue()) > 3 {
		if season, sep, episode, ok := extractSeasonAndEpisode(tkn.getValue()); ok {
			seasonPrefixTkn := newToken("S")
			seasonPrefixTkn.setIdentifiedKeywordCategory(keywordCatSeasonPrefix)
			seasonPrefixTkn.setKind(tokenKindCharacter)

			seasonTkn := newToken(season)
			seasonTkn.setMetadataKind(metadataKindSeason)
			seasonTkn.setKind(tokenKindNumber)

			sepTkn := newToken(sep)
			sepTkn.setIdentifiedKeywordCategory(keywordCatEpisodePrefix)
			sepTkn.setKind(tokenKindCharacter)

			episodeTkn := newToken(episode)
			episodeTkn.setMetadataKind(metadataKindEpisodeNumber)
			if isNumber(episode) {
				episodeTkn.setKind(tokenKindNumber)
			} else {
				episodeTkn.setKind(tokenKindNumberLike)
			}

			p.tokenManager.tokens.overwriteAndInsertManyAt(p.tokenManager.tokens.getIndexOf(tkn), []*token{seasonPrefixTkn, seasonTkn, sepTkn, episodeTkn})
			return true
		}
	}

	//if !ok {
	//	return []*token{}, false
	//}

	return false

}

// collectMetadata collects the metadata elements from the parsed tokens.
func (p *parser) collectMetadata() {

}
