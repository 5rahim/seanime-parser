package seanime_parser

import (
	"strconv"
	"strings"
)

func (p *parser) parseEpisode() {

	// Check alt episode number
	// e.g. 01 (12)
	for {
		found, tkns := p.tokenManager.tokens.findWithMetadataKind(metadataEpisodeNumber)
		if !found {
			break
		}

		last := tkns[len(tkns)-1]

		nextTkns, found, _ := p.tokenManager.tokens.getCategorySequenceAfter(p.tokenManager.tokens.getIndexOf(last), []tokenCategory{
			tokenCatOpeningBracket, // (
			tokenCatUnknown,        // 12
			tokenCatClosingBracket, // )
		}, true)
		if !found {
			break
		}

		if nextTkns[0].getValue() != "(" || !nextTkns[1].isNumberKind() || nextTkns[2].getValue() != ")" {
			break
		}

		// Update token
		nextTkns[1].setMetadataCategory(metadataEpisodeNumberAlt)
		break
	}

	// Check combined or separated keywords other than season prefixes
	// e.g. ED1, ED 1, OVAs 1-3, OVAs 1 ~ 3, OVA1, OVA 1v2
	p.parseKeywordsWithEpisodes()

	// Check last number before the first opening bracket (if there is one at the beginning [subgroup], then, before the second opening bracket)
	// e.g. Title - 01
	p.parseEpisodeBySearching(false)

	// e.g Title 01
	p.parseEpisodeBySearching(true)

}

// parseEpisodeBySearching parses episode numbers by searching for numbers that are not followed by a season prefix or episode prefix
// e.g. - 01, 1 [
func (p *parser) parseEpisodeBySearching(aggressive bool) {

	// Check if we already have an episode number
	found, _ := p.tokenManager.tokens.findWithMetadataKind(metadataEpisodeNumber)
	if found {
		return
	}

	// Check "- 01 [...]"
	for {
		var openingBracketTkn *token

		for _, tkn := range *p.tokenManager.tokens {
			if tkn.isOpeningBracket() {
				if p.tokenManager.tokens.getIndexOf(tkn) == 0 {
					continue
				}
				openingBracketTkn = tkn
				break
			}
		}

		if openingBracketTkn == nil {
			break
		}
		// Get previous token
		numTkn, found, _ := p.tokenManager.tokens.getTokenBeforeSD(openingBracketTkn)
		if !found {
			break
		}
		// Check if previous token is a number or number-like
		if !numTkn.isNumberOrLikeKind() || !numTkn.isUnknown() {
			break
		}
		if !aggressive && !p.tokenManager.tokens.foundDashSeparatorBefore(numTkn) {
			break

		}

		numTkn.setMetadataCategory(metadataEpisodeNumber)
		return // Found episode number, end

	}

	// Check for last number
	for {
		var lastNumTkn *token

		for _, tkn := range *p.tokenManager.tokens {
			if tkn.isYear() {
				continue
			}
			if tkn.isNumberOrLikeKind() && tkn.isUnknown() {
				lastNumTkn = tkn
			}
		}

		if lastNumTkn == nil {
			break
		}

		if !aggressive && !p.tokenManager.tokens.foundDashSeparatorBefore(lastNumTkn) {
			break
		}

		// When searching aggressively
		// Check that the number MIGHT be an episode number
		// e.g. if < 10 should be zero padded to avoid false positives with e.g. "Title 2"
		if aggressive {
			if isNumber(lastNumTkn.getValue()) {
				intVal, err := strconv.Atoi(lastNumTkn.getValue())
				if err != nil {
					break
				}
				if intVal < 10 && !isNumberZeroPadded(lastNumTkn.getValue()) {
					break
				}
			}
		}

		lastNumTkn.setMetadataCategory(metadataEpisodeNumber)
		return // Found episode number, end
	}

	//for _, tkn := range *p.tokenManager.tokens {
	//
	//	if tkn.isKeyword() || !tkn.isUnknown() { // Don't bother if token is already a keyword
	//		continue // Skip to next token
	//	}
	//
	//}
}

// parseKeywordsWithEpisodes parses keywords that are combined or separated with a number AND not a season prefix or episode prefix
// e.g. ED1, ED 1, OVAs 1-3, OVAs 1 ~ 3, OVA1, OVA 1v2
func (p *parser) parseKeywordsWithEpisodes() {

	for _, tkn := range *p.tokenManager.tokens {

		if tkn.isKeyword() || !tkn.isUnknown() { // Don't bother if token is already a keyword
			continue // Skip to next token
		}

		keywords, found := p.tokenManager.keywordManager.findKeywordsBy(func(kw *keyword) bool {
			// Get keywords that are combined or separated with a number AND not a season prefix or episode prefix
			// e.g. ED1, ED 1
			return (kw.isCombinedWithNumber() || kw.isSeparatedWithNumber()) &&
				!kw.isSeasonPrefix() && !kw.isEpisodePrefix() && !kw.isVolumePrefix() && !kw.isPartPrefix() && // Skip all these because they are handled in parseSeason()
				strings.HasPrefix(tkn.getNormalizedValue(), kw.value) // Token value starts with keyword value
		})

		if !found {
			continue // Skip to next token
		}

	keywordLoop:
		for _, keyword := range keywords {

			// e.g. ED1
			if keyword.isCombinedWithNumber() {

				// Check if prefix is followed by a number or number-like (e.g. 01, 01v2)
				remaining := strings.TrimPrefix(tkn.getNormalizedValue(), keyword.value)

				if len(remaining) > 0 && isNumberOrLike(remaining) {

					// e.g. ED
					prefixTkn := newToken(keyword.value)
					prefixTkn.setIdentifiedKeywordCategory(keyword.category)
					prefixTkn.setKind(tokenKindWord)

					// e.g. 01, 1, 3
					numberTkn := newToken(tkn.getValue()[len(keyword.value):])
					numberTkn.setMetadataCategory(metadataOtherEpisodeNumber)
					numberTkn.setKind(tokenKindNumberLike)
					if isNumber(remaining) {
						numberTkn.setKind(tokenKindNumber)
					}

					firstNumberIsZeroPadded := isNumberZeroPadded(remaining)

					// Overwrite token and insert new tokens
					// "ED1" -> "ED", "1"
					p.tokenManager.tokens.overwriteAndInsertManyAt(p.tokenManager.tokens.getIndexOf(tkn), []*token{prefixTkn, numberTkn})

					if isNumber(remaining) { // e.g. ED1.5, don't bother if ED1v2
						p.tokenManager.tokens.checkNumberWithDecimal(numberTkn) // Check if number is decimal
					}

					// Check range
					// e.g. ED1-3, ED01-03, ED1 ~ 3, ED01 ~ 03
					if nextNumTkn, found, kind := checkNumberRangeAfterToken(p, numberTkn, firstNumberIsZeroPadded); found {
						// e.g. ED1-3, ED01-03
						if kind == 0 {
							nextNumTkn.setMetadataCategory(metadataOtherEpisodeNumber)
							p.tokenManager.tokens.checkNumberWithDecimal(nextNumTkn) // Check if number is decimal
							break keywordLoop                                        // Skip to next token
						}
					}

					break keywordLoop // Skip to next token

				}

			}

			// e.g. ED 1
			if keyword.isSeparatedWithNumber() {

				// Get next token, by skipping delimiters
				// Check if next token is a number or number-like
				if nextTkn, found, _ := p.tokenManager.tokens.getTokenAfterSD(tkn); found &&
					(nextTkn.isNumberOrLikeKind() && nextTkn.isUnknown()) {

					tkn.setIdentifiedKeywordCategory(keyword.category)
					nextTkn.setMetadataCategory(metadataOtherEpisodeNumber)

					// Check range
					firstSeasonIsZeroPadded := isNumberZeroPadded(nextTkn.getValue())
					// e.g. ED 1-3, ED 01-03, ED 1 ~ 3, ED 01 ~ 03
					if nextNumTkn, found, kind := checkNumberRangeAfterToken(p, nextTkn, firstSeasonIsZeroPadded); found {
						// e.g. ED 1-3, ED 01-03
						if kind == 0 {
							nextNumTkn.setMetadataCategory(metadataOtherEpisodeNumber)
							break keywordLoop // Skip to next token
						}
					}

					break keywordLoop // Skip to next token

				}

			}

		}
	}

}
