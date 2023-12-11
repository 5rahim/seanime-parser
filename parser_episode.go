package seanime_parser

import (
	"strconv"
	"strings"
)

func (p *parser) parseEpisode() {

	foundEpisode := false

	// Check alt episode number or range //TODO
	// e.g. 01 (12)
	for {
		found, tkns := p.tokenManager.tokens.findWithMetadataCategory(metadataEpisodeNumber)
		if !found {
			break
		}

		foundEpisode = true

		last := tkns[len(tkns)-1]

		nextTkns, found, _ := p.tokenManager.tokens.getCategorySequenceAfter(p.tokenManager.tokens.getIndexOf(last), []tokenCategory{
			tokenCatOpeningBracket, // (
			tokenCatUnknown,        // 12
			tokenCatClosingBracket, // )
		}, true)
		if !found {

			{ // Check range after found episode
				if len(tkns) != 1 {
					break
				}
				// Check range
				rangeTkns, foundRange := p.tokenManager.tokens.checkNumberRangeAfter(last)
				if !foundRange {
					break
				}
				rangeTkns[1].setMetadataCategory(metadataEpisodeNumber)
			}
			break
		}

		if nextTkns[0].getValue() != "(" || !nextTkns[1].isNumberKind() || nextTkns[2].getValue() != ")" {
			break
		}

		// Update token
		nextTkns[1].setMetadataCategory(metadataEpisodeNumberAlt)
		break
	}

	// Search by alt episode number
	// We check if any unknown number token is followed by "({number})", e.g. {tkn} (14)
	for _, numTkn := range *p.tokenManager.tokens {
		if !numTkn.isNumberOrLikeKind() || !numTkn.isUnknown() {
			continue // Check next token
		}
		nextTkns, found, _ := p.tokenManager.tokens.getCategorySequenceAfter(p.tokenManager.tokens.getIndexOf(numTkn), []tokenCategory{
			tokenCatOpeningBracket, // (
			tokenCatUnknown,        // 12
			tokenCatClosingBracket, // )
		}, true)
		if !found {
			continue // Check next token
		}
		if nextTkns[0].getValue() != "(" || !nextTkns[1].isNumberKind() || nextTkns[2].getValue() != ")" {
			continue
		}

		// Update tokens
		numTkn.setMetadataCategory(metadataEpisodeNumber)
		nextTkns[1].setMetadataCategory(metadataEpisodeNumberAlt)
		foundEpisode = true
		break
	}

	if foundEpisode {
		return
	}

	// Check combined or separated keywords other than season prefixes
	// e.g. ED1, ED 1, OVAs 1-3, OVAs 1 ~ 3, OVA1, OVA 1v2
	if foundEpisode := p.parseKeywordsWithEpisodes(); foundEpisode {
		return // Stop if we find the actual episode number (e.g. Ep01) (not keyword number like OVA 1)
	}

	// e.g. 01 of 24
	if found := p.parseEpisodeByRangeSeparator("OF"); found {
		return
	}

	// Check last number before the first opening bracket (if there is one at the beginning [subgroup], then, before the second opening bracket)
	// e.g. Title - 01
	if found := p.parseEpisodeBySearching(false); found {
		return
	}

	// e.g. [12]
	if found := p.parseEpisodeByEnclosedNumber(); found {
		return
	}

	// e.g Title 01
	if found := p.parseEpisodeBySearching(true); found {
		return
	}

}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// parseEpisodeBySearching parses the episode number by searching for different patterns.
func (p *parser) parseEpisodeBySearching(aggressive bool) bool {

	// Check "- 01 [...]" or "- 01 480p"
	for {
		var openingBracketTkn *token

		for _, tkn := range *p.tokenManager.tokens {
			// Find the first opening bracket or metadata keyword, whatever comes first
			if tkn.isOpeningBracket() || tkn.isKeyword() {
				if p.tokenManager.tokens.getIndexOf(tkn) == 0 { // Skip first token if it's an opening bracket
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
		// Check if previous token is an unknown number or number-like
		if !numTkn.isNumberOrLikeKind() || !numTkn.isUnknown() {
			break
		}

		// Check we find a range before
		// e.g. "1 - {numTkn} [...]"
		if rangeTkns, found := p.tokenManager.tokens.checkNumberRangeBefore(numTkn); found {
			rangeTkns[1].setMetadataCategory(metadataEpisodeNumber)
			numTkn.setMetadataCategory(metadataEpisodeNumber)
			return true // Found episode number, end
		}

		// If we are not searching aggressively, check if there is a dash separator before the number
		// e.g. "- 01"
		if !aggressive && !p.tokenManager.tokens.foundDashSeparatorBefore(numTkn) {
			break
		}

		// When searching aggressively
		// Check that the number might really be an episode number
		// e.g. if {lastNumTkn} < 10, lastNumTkn should be zero padded to avoid false positives like "Title 2"
		//FIXME False positive: "Evangelion 3.0 You Can (Not) Redo" -> "3.0" is identified as episode number
		if aggressive {
			if numTkn.isNumberKind() {
				intVal, err := strconv.Atoi(numTkn.getValue())
				if err != nil {
					break
				}
				println(intVal)
				if intVal < 10 && !isNumberZeroPadded(numTkn.getValue()) {
					break
				}
				// should be isolated
				if !p.tokenManager.tokens.isIsolated(numTkn) {
					break
				}
			}
			// in the case of "Title 2v2", we can safely identify "2v2" as an episode number
		}

		numTkn.setMetadataCategory(metadataEpisodeNumber)
		return true // Found episode number, end

	}

	// Check for last unknown number
	for {
		var lastNumTkn *token

		// Get the last unknown number token
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

		// Check we find a range before
		// e.g. "1 - {lastNumTkn} [...]"
		if rangeTkns, found := p.tokenManager.tokens.checkNumberRangeBefore(lastNumTkn); found {
			rangeTkns[1].setMetadataCategory(metadataEpisodeNumber)
			lastNumTkn.setMetadataCategory(metadataEpisodeNumber)
			return true // Found episode number, end
		}

		// If we are not searching aggressively, check if there is a dash separator before the number
		// e.g. "- 01"
		if !aggressive && !p.tokenManager.tokens.foundDashSeparatorBefore(lastNumTkn) {
			break
		}

		// When searching aggressively
		// Check that the number might really be an episode number
		// e.g. if {lastNumTkn} < 10, lastNumTkn should be zero padded to avoid false positives like "Title 2"
		//FIXME False positive: "Evangelion 3.0 You Can (Not) Redo" -> "3.0" is identified as episode number
		if aggressive {
			if lastNumTkn.isNumberKind() {
				intVal, err := strconv.Atoi(lastNumTkn.getValue())
				if err != nil {
					break
				}
				println(intVal)
				if intVal < 10 && !isNumberZeroPadded(lastNumTkn.getValue()) {
					break
				}
				// should be isolated
				if !p.tokenManager.tokens.isIsolated(lastNumTkn) {
					break
				}
			}
			// in the case of "Title 2v2", we can safely identify "2v2" as an episode number
		}

		lastNumTkn.setMetadataCategory(metadataEpisodeNumber)
		return true // Found episode number, end
	}

	return false
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// parseKeywordsWithEpisodes parses keywords that are combined or separated with a number.
// It does not handle season/volume/part prefixes (keywordCatSeasonPrefix, ...) as those are handled by parseSeason.
//
// It handles the following cases:
// keywordKindCombinedWithNumber
// keywordKindSeparatedWithNumber
//
// e.g. Ep1, ED1, ED 1, OVAs 1-3, OVAs 1 ~ 3, OVA1, OVA 1v2
func (p *parser) parseKeywordsWithEpisodes() (foundEpisode bool) {

	for _, tkn := range *p.tokenManager.tokens {

		if tkn.isKeyword() || !tkn.isUnknown() { // Don't bother if token is already a keyword
			continue // Skip to next token
		}

		keywords, found := p.tokenManager.keywordManager.findKeywordsBy(func(kw *keyword) bool {
			// Get keywords that are combined or separated with a number AND not a season prefix or episode prefix
			// e.g. ED1, ED 1
			return (kw.isCombinedWithNumber() || kw.isSeparatedWithNumber()) &&
				!kw.isSeasonPrefix() && !kw.isVolumePrefix() && !kw.isPartPrefix() && // Skip all these because they are handled in parseSeason()
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
					if keyword.isEpisodePrefix() {
						numberTkn.setMetadataCategory(metadataEpisodeNumber)
					}
					numberTkn.setKind(tokenKindNumberLike)
					if isNumber(remaining) {
						numberTkn.setKind(tokenKindNumber)
					}

					firstNumberIsZeroPadded := isNumberZeroPadded(remaining)

					// Overwrite token and insert new tokens
					// "ED1" -> "ED", "1"
					p.tokenManager.tokens.overwriteAndInsertManyAt(p.tokenManager.tokens.getIndexOf(tkn), []*token{prefixTkn, numberTkn})
					foundEpisode = true

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
					if keyword.isEpisodePrefix() {
						foundEpisode = true
						nextTkn.setMetadataCategory(metadataEpisodeNumber)
					}

					// Check range
					firstSeasonIsZeroPadded := isNumberZeroPadded(nextTkn.getValue())
					// e.g. ED 1-3, ED 01-03, ED 1 ~ 3, ED 01 ~ 03
					if nextNumTkn, found, kind := checkNumberRangeAfterToken(p, nextTkn, firstSeasonIsZeroPadded); found {
						// e.g. ED 1-3, ED 01-03
						if kind == 0 {
							nextNumTkn.setMetadataCategory(metadataOtherEpisodeNumber)
							if keyword.isEpisodePrefix() {
								nextNumTkn.setMetadataCategory(metadataEpisodeNumber)
							}
							break keywordLoop // Skip to next token
						}
					}

					break keywordLoop // Skip to next token

				}

			}

		}
	}

	return

}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (p *parser) parseEpisodeByEnclosedNumber() bool {

	if !p.tokenManager.tokens.allUnknownTokensAreEnclosed() {
		return false
	}

	for _, numTkn := range *p.tokenManager.tokens {
		if !numTkn.isUnknown() || !numTkn.isNumberOrLikeKind() || !numTkn.isEnclosed() {
			continue // Check next token
		}

		// e.g. [
		prevTkn, found := p.tokenManager.tokens.getTokenBefore(numTkn)
		if !found || !prevTkn.isOpeningBracket() {
			continue
		}
		// e.g. ]
		nextTkn, found := p.tokenManager.tokens.getTokenAfter(numTkn)
		if !found || !nextTkn.isClosingBracket() {
			continue
		}

		// e.g. We found [12] or (12)
		numTkn.setMetadataCategory(metadataEpisodeNumber)
		return true // Found

	}

	return false
}

func (p *parser) parseEpisodeByRangeSeparator(value string) bool {

	for _, numTkn := range *p.tokenManager.tokens {
		if !numTkn.isUnknown() || !numTkn.isNumberOrLikeKind() {
			continue // Check next token
		}

		// e.g. of
		ofTkn, found, _ := p.tokenManager.tokens.getTokenAfterSD(numTkn)
		if !found || ofTkn.getNormalizedValue() != value {
			continue
		}
		// e.g. [
		secondNumTkn, found, _ := p.tokenManager.tokens.getTokenAfterSD(ofTkn)
		if !found || !secondNumTkn.isNumberOrLikeKind() {
			continue
		}

		// e.g. We found "01 of 24"
		numTkn.setMetadataCategory(metadataEpisodeNumber)
		ofTkn.setCategory(tokenCatKnown)
		secondNumTkn.setMetadataCategory(metadataOtherEpisodeNumber)
		return true

	}
	return false

}
