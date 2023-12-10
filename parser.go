package seanime_parser

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

func Parse(filename string) *Metadata {

	p := newParser(filename)
	p.parse()
	p.cleanUp()

	return p.metadata

}

func (p *parser) parse() {

	p.parseKeywords("normal")

	p.parseSeason()

	p.parseEpisode()

	p.parseEpisodeTitle()

	p.parseTitle()

	p.parseReleaseGroup()

	p.parseKeywords("")

	p.collectMetadata()

}

func (p *parser) parseKeywords(priority string) {

	for _, tkn := range *p.tokenManager.tokens {

		if tkn.isKeyword() || !tkn.isUnknown() { // Don't bother if token is already a keyword
			continue // Skip to next token
		}

		// Identify keyword
		_ = p.identifyKeyword(tkn, priority)

	}

}

// identifyKeyword identifies STANDALONE and multi-PART keywords for the given token
func (p *parser) identifyKeyword(tkn *token, priority string) bool {

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

	// Check if token is a known pre-defined STANDALONE keyword (e.g. "60FPS")
	if len(tkn.getValue()) > 1 {
		if keyword, found := p.tokenManager.keywordManager.findStandaloneKeywordByValue(tkn.getValue()); found {

			// When the priority is "normal", we only want to identify STANDALONE keywords that are not an anime type
			// That is because those are prone to false positives
			if priority == "normal" && keyword.isAnimeType() {
				return false
			}

			// when the priority is "normal", we only want to identify STANDALONE keywords that are not ambiguous
			if priority == "normal" && p.tokenManager.keywordManager.isKeywordAmbiguous(keyword) {
				return false
			}

			tkn.setIdentifiedKeywordCategory(keyword.category)
			return true
		}
	}

	return false

}

// collectMetadata collects the metadata elements from the parsed tokens.
// de-duplicates elements
func (p *parser) collectMetadata() {

	for _, tkn := range *p.tokenManager.tokens {

		if !tkn.isIdentifiedMetadata() {
			continue
		}

	}

}

func (p *parser) cleanUp() {

	// TODO: Remove versions from numbers (e.g. 01v2 -> 01) and add ReleaseVersion to metadata

}
