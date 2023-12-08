package seanime_parser

import (
	"github.com/davecgh/go-spew/spew"
	"strconv"
	"strings"
)

func (p *parser) parseSeason() {

	for _, tkn := range *p.tokenManager.tokens {

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
				continue // Skip to next token
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
					// /!\ This section will hydrate the METADATA
					if keyword.isCombinedWithNumber() {

						// Check if token is after file metadata
						if p.tokenManager.tokens.isTokenAfterFileMetadata(tkn) {
							continue // Skip to next token
						}

						// Check if prefix is followed by a number or number-like (e.g. 01, 01v2)
						remaining := strings.TrimPrefix(tkn.getNormalizedValue(), keyword.value)

						if len(remaining) > 0 && isNumberOrLike(remaining) {

							// e.g. S
							seasonPrefixTkn := newToken(keyword.value)
							seasonPrefixTkn.setIdentifiedKeywordCategory(keywordCatSeasonPrefix)
							seasonPrefixTkn.setKind(tokenKindWord)

							// e.g. 01, 1, 3
							seasonTkn := newToken(tkn.getValue()[len(keyword.value):])
							seasonTkn.setMetadataCategory(metadataSeason)
							seasonTkn.setKind(tokenKindNumberLike)
							if isNumber(remaining) {
								seasonTkn.setKind(tokenKindNumber)
							}

							firstSeasonIsZeroPadded := isNumberZeroPadded(remaining)

							p.tokenManager.tokens.overwriteAndInsertManyAt(p.tokenManager.tokens.getIndexOf(tkn), []*token{seasonPrefixTkn, seasonTkn})

							spew.Dump(seasonTkn)

							// Check range
							// e.g. S1-3, S01-03, S1 ~ 3, S01 ~ 03
							if rangeTkns, found, dlSkipped := p.tokenManager.tokens.getCategorySequenceAfter(p.tokenManager.tokens.getIndexOf(seasonTkn), []tokenCategory{
								tokenCatSeparator, // -
								tokenCatUnknown,   // 05
							}, true); found {

								// Check episode
								if rangeTkns[1].isNumberOrLikeKind() && rangeTkns[0].isDashSeparator() {

									// e.g. S1 - 03, S01- 03
									if dlSkipped > 0 {
										if intVal, err := strconv.Atoi(rangeTkns[1].getValue()); err == nil {
											// e.g. if < 10 -> 01, 02, 03. if > 10 -> 11, 12, 13
											if (intVal < 10 && isNumberZeroPadded(rangeTkns[1].getValue())) || (intVal > 10) {
												rangeTkns[1].setMetadataCategory(metadataEpisodeNumber)
												continue // Skip to next token
											}
										} else { // /!\ might need to do some additional checks on the number
											rangeTkns[1].setMetadataCategory(metadataEpisodeNumber)
											continue // Skip to next token
										}

										// e.g. S1-03 (dlSkipped = 0) Where 03 might be an episode. This is not very likely
									} else if !firstSeasonIsZeroPadded && isNumberZeroPadded(rangeTkns[1].getValue()) {
										rangeTkns[1].setMetadataCategory(metadataSeason)
										continue // Skip to next token
									}
								}

								if !rangeTkns[1].isNumberKind() {
									continue // Skip to next token
								}

								intVal, err := strconv.Atoi(rangeTkns[1].getValue())
								if err != nil {
									continue // Skip to next token
								}

								if intVal < 10 && (firstSeasonIsZeroPadded && !isNumberZeroPadded(rangeTkns[1].getValue())) {
									continue // Skip to next token
								}

								// e.g. S1-3
								rangeTkns[1].setMetadataCategory(metadataSeason)

							}

							continue // Skip to next token
						}

					}

					// e.g. Season 01
					if keyword.isSeparatedWithNumber() {

						// Get next token, by skipping delimiters
						// Check if next token is a number or number-like
						if nextTkn, found, _ := p.tokenManager.tokens.getTokenAfterSD(tkn); found &&
							(nextTkn.isNumberOrLikeKind()) {

							nextTkn.setMetadataCategory(metadataSeason)
							continue // Skip to next token

						}

					}

					// e.g. 1st Season, first season
					if keyword.isOrdinalSuffix() {

						// Get previous token, by skipping delimiters
						// Check if next token is an ordinal number
						if nextTkn, found, _ := p.tokenManager.tokens.getTokenBeforeSD(tkn); found &&
							(nextTkn.isOrdinalNumber()) {

							if num, ok := getNumberFromOrdinal(nextTkn.getValue()); ok {
								nextTkn.setValue(strconv.Itoa(num))
								nextTkn.setMetadataCategory(metadataSeason)
								nextTkn.setKind(tokenKindNumber)
								continue // Skip to next token
							}

						}

					}

				}
			}

		}

		// Parse 01x01
		if strings.Contains(tkn.getNormalizedValue(), "X") && len(tkn.getValue()) > 3 {
			// Extract season and episode
			if season, sep, episode, ok := extractSeasonAndEpisode(tkn.getValue()); ok {
				if len(season) > 2 {
					continue // Skip to next token
				}

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
				continue // Skip to next token
			}
		}

	}

}

func checkNumberRangeAfterToken(p *parser, tkn *token, prevNumberIsPadded bool) (*token, bool, int) {

	var nextNumTkn *token
	found := false
	var kind int // 0 = season, 1 = episode

	// Check range
	// e.g. S1-3, S01-03, S1 ~ 3, S01 ~ 03
	if rangeTkns, ok, dlSkipped := p.tokenManager.tokens.getCategorySequenceAfter(p.tokenManager.tokens.getIndexOf(tkn), []tokenCategory{
		tokenCatSeparator, // -
		tokenCatUnknown,   // 05
	}, true); ok {

		// Check episode
		if rangeTkns[1].isNumberOrLikeKind() && rangeTkns[0].isDashSeparator() {

			nextNumTkn = rangeTkns[1]

			// e.g. S1 - 03, S01- 03
			if dlSkipped > 0 {
				if intVal, err := strconv.Atoi(rangeTkns[1].getValue()); err == nil {
					// e.g. if < 10 -> 01, 02, 03. if > 10 -> 11, 12, 13
					if (intVal < 10 && isNumberZeroPadded(rangeTkns[1].getValue())) || (intVal > 10) {
						rangeTkns[1].setMetadataCategory(metadataEpisodeNumber)
						kind = 1
						found = true
					}
				} else { // /!\ might need to do some additional checks on the number
					rangeTkns[1].setMetadataCategory(metadataEpisodeNumber)
					kind = 1
					found = true
				}

				// e.g. S1-03 (dlSkipped = 0) Where 03 might be an episode. This is not very likely
			} else if !prevNumberIsPadded && isNumberZeroPadded(rangeTkns[1].getValue()) {
				rangeTkns[1].setMetadataCategory(metadataSeason)
				kind = 1
				found = true
			}
		}

		// Avoid this case: S1-2v2
		if !rangeTkns[1].isNumberKind() {
			found = false
		}

		intVal, err := strconv.Atoi(rangeTkns[1].getValue())
		if err != nil {
			found = false
		}

		// Avoid this case: S01 - 3
		if intVal < 10 && (prevNumberIsPadded && !isNumberZeroPadded(rangeTkns[1].getValue())) {
			found = false
		}

		// e.g. S1-3
		rangeTkns[1].setMetadataCategory(metadataSeason)

	}

	return nextNumTkn, found, kind

}
