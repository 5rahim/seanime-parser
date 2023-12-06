package seanime_parser

func (t *token) getValue() string {
	if t == nil {
		return ""
	}
	return t.Value
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
