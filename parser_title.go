package seanime_parser

func (p *parser) parseEpisodeTitle() {

	// Get all tokens after the last episode number token and before an opening bracket/file info metadata/EOF
	found, epTkns := p.tokenManager.tokens.findWithMetadataCategory(metadataEpisodeNumber)
	if !found {
		return
	}
	lastEpTkn := epTkns[0]
	if len(epTkns) > 0 {
		lastEpTkn = epTkns[len(epTkns)-1]
	}

	tkns, found := p.tokenManager.tokens.walkAndCollecIf(
		p.tokenManager.tokens.getIndexOf(lastEpTkn)+1,
		func(tkn *token) bool {
			return tkn.isUnknown() && !tkn.isKeyword() && !tkn.isSeparator()
		},
		func(tkn *token) bool {
			return (tkn.isOpeningBracket() && tkn.getValue() == "[") || tkn.isKeyword()
		})
	if !found {
		return
	}

	p.tokenManager.tokens.combineTitle(tkns[0], tkns[len(tkns)-1], metadataEpisodeTitle)

	return

}

func (p *parser) parseTitleIfAllEnclosed() (foundTitle bool) {
	foundTitle = false

	if !p.tokenManager.tokens.allUnknownTokensAreEnclosed() {
		return
	}

	// Get enclosed tokens before enclosed episode number
	// e.g. "[sub][anime title][01][...]"
	for {
		// Get first episode number token
		found, epTkns := p.tokenManager.tokens.findWithMetadataCategory(metadataEpisodeNumber)
		if !found {
			break // Next try
		}
		firstEpTkn := epTkns[0]

		// Get the first and second opening bracket going backwards
		// e.g. `[` <- anime title ] `[` <- 01 ][ ... ]
		firstOpeningBracketTkn, found := p.tokenManager.tokens.getFirstOccurrenceBefore(
			p.tokenManager.tokens.getIndexOf(firstEpTkn),
			func(tkn *token) bool {
				return tkn.isOpeningBracket()
			})
		if !found {
			break // Next try
		}
		// Get second opening bracket going backwards
		secondOpeningBracketTkn, found := p.tokenManager.tokens.getFirstOccurrenceBefore(
			p.tokenManager.tokens.getIndexOf(firstOpeningBracketTkn),
			func(tkn *token) bool {
				return tkn.isOpeningBracket()
			})
		if !found {
			break // Next try
		}
		// Get all unknown tokens between the two opening brackets
		// e.g. `[` -> "anime" -> "title" ] `[` 01 ][ ... ]
		tkns, found := p.tokenManager.tokens.walkAndCollecIf(
			p.tokenManager.tokens.getIndexOf(secondOpeningBracketTkn)+1,
			func(tkn *token) bool {
				return tkn.isUnknown() && !tkn.isKeyword() && !tkn.isSeparator()
			},
			func(tkn *token) bool {
				// Stop when we encounter the first opening bracket or a keyword
				// e.g. [Mobile_Suit_Gundam_Seed_Destiny_HD_REMASTER][07] -> Mobile Suit Gundam Seed Destiny
				return tkn.UUID == firstOpeningBracketTkn.UUID || tkn.isKeyword()
			})
		if !found {
			break // Next try
		}

		p.tokenManager.tokens.combineTitle(tkns[0], tkns[len(tkns)-1], metadataTitle)
		return true
	}

	// Get the second enclosed group
	// e.g. "[sub][anime title][BDRIP][...]"
	for {
		// Get all opening brackets
		openingBracketTkns, found := p.tokenManager.tokens.filter(func(tkn *token) bool {
			return tkn.isOpeningBracket() && tkn.getValue() != "("
		})
		if !found {
			break // Next try
		}

		if len(openingBracketTkns) < 2 {
			break
		}

		// Get all unknown tokens between the second and third opening brackets
		// e.g. [ sub ]`[` -> "anime" -> "title" -> `]` [ ... ]
		tkns, found := p.tokenManager.tokens.walkAndCollecIf(
			p.tokenManager.tokens.getIndexOf(openingBracketTkns[1])+1,
			func(tkn *token) bool {
				return tkn.isUnknown() && !tkn.isKeyword() && !tkn.isSeparator()
			},
			func(tkn *token) bool {
				// Stop when we encounter the third opening bracket or a keyword
				return tkn.UUID == openingBracketTkns[2].UUID || tkn.isKeyword()
			})

		if !found {
			break // Next try
		}

		p.tokenManager.tokens.combineTitle(tkns[0], tkns[len(tkns)-1], metadataTitle)
		return true
	}

	return
}

func (p *parser) parseTitle() {

	if found := p.parseTitleIfAllEnclosed(); found {
		return
	}

	// e.g. [sub] anime title [...]
	for {
		// Get first non-enclosed token
		nonEnclosedTkns, found := p.tokenManager.tokens.filter(func(tkn *token) bool {
			return !tkn.isEnclosed() && tkn.isUnknown() && !tkn.isKeyword()
		})
		if !found {
			break // Next try
		}

		// Collect all unknown tokens from the first non-enclosed token until an opening bracket or keyword is found
		// e.g. "[ignored] collected collected [ignored]"
		tkns, found := p.tokenManager.tokens.walkAndCollecIf(
			p.tokenManager.tokens.getIndexOf(nonEnclosedTkns[0]),
			func(tkn *token) bool {
				return !tkn.isEnclosed() && // not enclosed
					tkn.isUnknown() && // unknown
					!tkn.isKeyword() // not a keyword
				//!tkn.isSeparator() // not a separator
			},
			func(tkn *token) bool {
				return (tkn.isOpeningBracket() && tkn.getValue() == "[") || tkn.isKeyword()
			})
		if !found {
			break // Next try
		}
		// Title should not be after file info metadata like 1080p
		if p.tokenManager.tokens.isTokenAfterFileMetadata(tkns[0]) {
			break // Next try
		}

		p.tokenManager.tokens.combineTitle(tkns[0], tkns[len(tkns)-1], metadataTitle)
		return
	}

}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (p *parser) parseReleaseGroup() {

	foundTitle, titleTkns := p.tokenManager.tokens.findWithMetadataCategory(metadataTitle)

	// Handle case where all unknown tokens are enclosed
	for {
		if !foundTitle {
			break // Next try
		}

		titleTkn := titleTkns[0]

		// Get all unknown tokens before the title
		unknownTkns, found := p.tokenManager.tokens.walkAndCollecIf(
			0,
			func(tkn *token) bool {
				return tkn.isUnknown() && !tkn.isKeyword() && !tkn.isSeparator() && tkn.isEnclosed()
			},
			func(tkn *token) bool {
				return tkn.UUID == titleTkn.UUID
			})
		if !found {
			break // Next try
		}

		// Found release group
		if len(unknownTkns) == 1 {
			unknownTkns[0].setMetadataCategory(metadataReleaseGroup)
			return
		}

		// Found longer release group
		p.tokenManager.tokens.combineTitle(unknownTkns[0], unknownTkns[len(unknownTkns)-1], metadataReleaseGroup)
		return
	}

	// If we still haven't found a release group, try to find:
	// - the first enclosed group of unknown tokens going backwards
	for {

		// Get all closing brackets
		closingBracketTkns, found := p.tokenManager.tokens.filter(func(tkn *token) bool {
			return tkn.isClosingBracket() && tkn.getValue() != ")"
		})
		if !found || len(closingBracketTkns) == 1 {
			break // Next try
		}
		lastClosingBracket := closingBracketTkns[len(closingBracketTkns)-1]

		lastOpeningBracket, found := p.tokenManager.tokens.getFirstOccurrenceBefore(
			p.tokenManager.tokens.getIndexOf(lastClosingBracket),
			func(tkn *token) bool {
				return tkn.isOpeningBracket() && isMatchingClosingBracket(tkn.getValue(), lastClosingBracket.getValue())
			})
		if !found {
			break // Next try
		}

		// Get all tokens between the opening and closing brackets
		unknownTkns, found := p.tokenManager.tokens.walkAndCollecIf(
			p.tokenManager.tokens.getIndexOf(lastOpeningBracket)+1,
			func(tkn *token) bool {
				return tkn.isUnknown() && !tkn.isKeyword() && !tkn.isSeparator()
			},
			func(tkn *token) bool {
				return tkn.UUID == lastClosingBracket.UUID
			})
		if !found {
			break // Next try
		}

		// Found release group
		if len(unknownTkns) == 1 {
			unknownTkns[0].setMetadataCategory(metadataReleaseGroup)
			return
		}

		// Found longer release group
		p.tokenManager.tokens.combineTitle(unknownTkns[0], unknownTkns[len(unknownTkns)-1], metadataReleaseGroup)
		return
	}

	for {
		// Get all unknown tokens
		unknownTkns, found := p.tokenManager.tokens.filter(func(tkn *token) bool {
			return tkn.isUnknown() && !tkn.isKeyword() && !tkn.isSeparator()
		})
		if !found {
			break // Next try
		}

		unknownTkn := unknownTkns[0]
		if len(unknownTkns) > 1 {
			unknownTkn = unknownTkns[len(unknownTkns)-1]
		}

		if p.tokenManager.tokens.isTokenInFirstHalf(unknownTkn) {
			break
		}

		// Found release group
		unknownTkn.setMetadataCategory(metadataReleaseGroup)
		return
	}

}
