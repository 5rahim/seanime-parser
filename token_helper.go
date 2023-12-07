package seanime_parser

const tokenIdentifiedAsParts = 0
const tokenIdentifiedAsKeyword = 1

func (t *token) isIdentified() (int, bool) {
	if t == nil {
		return -1, false
	}
	if t.IdentifiedKeywordCategory != keywordCatNone {
		return tokenIdentifiedAsKeyword, true
	}
	if t.Parts != nil {
		return tokenIdentifiedAsParts, true
	}
	return -1, false
}
func (t *token) getValue() string {
	if t == nil {
		return ""
	}
	return t.Value
}

func (t *token) getCategory() tokenCategory {
	if t == nil {
		return ""
	}
	return t.Category
}

func (t *token) getIdentifiedKeywordCategory() (keywordCategory, bool) {
	if t == nil {
		return 0, false
	}
	return t.IdentifiedKeywordCategory, t.IdentifiedKeywordCategory != keywordCatNone
}

func (t *token) getKind() tokenKind {
	if t == nil {
		return ""
	}
	return t.Kind
}

func (t *token) isCategory(cat tokenCategory) bool {
	if t == nil {
		return false
	}
	return t.Category == cat
}

func (t *token) isKeyword() bool {
	if t == nil {
		return false
	}
	return t.IdentifiedKeywordCategory != keywordCatNone
}

func (t *token) isKeywordCategory(c keywordCategory) bool {
	if t == nil {
		return false
	}
	return t.IdentifiedKeywordCategory == c
}

func (t *token) isKind(kind tokenKind) bool {
	if t == nil {
		return false
	}
	return t.Kind == kind
}

func (t *token) isOpeningBracket() bool {
	if t == nil {
		return false
	}
	return t.Category == tokenCatOpeningBracket
}

func (t *token) isClosingBracket() bool {
	if t == nil {
		return false
	}
	return t.Category == tokenCatClosingBracket
}

func (t *token) isUnknown() bool {
	if t == nil {
		return false
	}
	return t.Category == tokenCatUnknown
}

func (t *token) isDelimiter() bool {
	if t == nil {
		return false
	}
	return t.Category == tokenCatDelimiter
}

func (t *token) isSeparator() bool {
	if t == nil {
		return false
	}
	return t.Category == tokenCatSeparator
}

func (t *token) isDotSeparator() bool {
	if t == nil {
		return false
	}
	return t.Category == tokenCatSeparator && t.Value == "."
}
func (t *token) isPlusSeparator() bool {
	if t == nil {
		return false
	}
	return t.Category == tokenCatSeparator && t.Value == "+"
}

func (t *token) isDashSeparator() bool {
	if t == nil {
		return false
	}
	return t.Category == tokenCatSeparator && t.Value == "-"
}

func (t *token) isEnclosed() bool {
	if t == nil {
		return false
	}
	return t.Enclosed
}

// isFileInfoMetadata checks if the token is classified as file information metadata.
// It returns true if the token DOES NOT belong to one of the following identified keyword categories:
//
// - keywordCatNone
// - keywordCatYear
// - keywordCatReleaseVersion
// - keywordCatReleaseGroup
// - keywordCatSeasonPrefix
// - keywordCatPartPrefix
// - keywordCatVolumePrefix
// - keywordCatEpisodePrefix
func (t *token) isFileInfoMetadata() bool {
	if t == nil {
		return false
	}

	switch t.IdentifiedKeywordCategory {
	case keywordCatNone:
		return false
	case keywordCatYear:
		return false
	case keywordCatReleaseVersion:
		return false
	case keywordCatReleaseGroup:
		return false
	case keywordCatSeasonPrefix:
		return false
	case keywordCatPartPrefix:
		return false
	case keywordCatVolumePrefix:
		return false
	case keywordCatEpisodePrefix:
		return false
	default:
		return true
	}
}

func (t *token) isAnimeInfoMetadata() bool {
	if t == nil {
		return false
	}
	return !t.isFileInfoMetadata()
}
