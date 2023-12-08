package seanime_parser

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
	// -> Check range

	// Check last number before the first opening bracket (if there is one at the beginning [subgroup], then, before the second opening bracket)
	// -> Check range

}
