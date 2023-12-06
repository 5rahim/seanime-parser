package seanime_parser

import (
	"github.com/google/uuid"
)

type tokenCategory = string

const (
	tokenCatUnknown        tokenCategory = "unknown"
	tokenCatDelimiter      tokenCategory = "delimiter"
	tokenCatSeparator      tokenCategory = "separator"
	tokenCatKnown          tokenCategory = "known"
	tokenCatParts          tokenCategory = "parts"
	tokenCatOpeningBracket tokenCategory = "openingBracket"
	tokenCatClosingBracket tokenCategory = "closingBracket"
)

type tokenKind = string

const (
	tokenKindUnknown          tokenKind = "unknown"
	tokenKindCharacter        tokenKind = "character"
	tokenKindWord             tokenKind = "word"
	tokenKindNumber           tokenKind = "number"
	tokenKindNumberLike       tokenKind = "numberLike"
	tokenKindOrdinalNumber    tokenKind = "ordinalNumber"
	tokenKindCrc32            tokenKind = "crc32"
	tokenKindPossibleVideoRes tokenKind = "possibleVideoRes"
	tokenKindYear             tokenKind = "year"
)

type token struct {
	UUID         string
	Value        string
	Category     tokenCategory
	Kind         tokenKind
	Enclosed     bool
	Parts        []*token
	MetadataKind metadataKind
}

func newToken(value string) *token {
	return &token{
		UUID:         uuid.NewString(),
		Value:        value,
		Category:     tokenCatUnknown,
		Kind:         tokenKindUnknown,
		Enclosed:     false,
		Parts:        nil,
		MetadataKind: 0,
	}
}

// setMetadataKind will update the token's MetadataKind
// and update its Category to tokenCatKnown if the metadataKind is not metadataKindUnknown.
func (t *token) setMetadataKind(mk metadataKind) {
	if t == nil {
		return
	}

	t.MetadataKind = mk
	t.Category = tokenCatKnown
	if mk == metadataKindUnknown {
		t.Category = tokenCatUnknown
	}
}

func (t *token) setCategory(c tokenCategory) {
	if t == nil {
		return
	}

	t.Category = c
}

func (t *token) setKind(k tokenKind) {
	if t == nil {
		return
	}

	t.Kind = k
}

func (t *token) setParts(p []*token) {
	if t == nil {
		return
	}

	t.Category = tokenCatParts
	t.Parts = p
}

func (t *token) setEnclosed(v bool) {
	if t == nil {
		return
	}

	t.Enclosed = v
}

func (t *token) getParts() (tokenParts []*token, found bool) {
	if t.Parts == nil {
		found = false
	}
	tokenParts = t.Parts
	found = true
	return
}
