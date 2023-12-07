package seanime_parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeywordGroups(t *testing.T) {

	tests := []struct {
		name             string
		input            string
		expectedTknValue string
	}{
		{"1", "BLU RAY 1080P", "BLU RAY"},
		{"2", "BLU-RAY 1080P", "BLU-RAY"},
		{"3", "TV RIP 1080P", "TV RIP"},
		{"4", "10 bits 1080P", "10 bits"},
		{"4", "10-bit 1080P", "10-bit"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := newTokenManager(tt.input)
			tkn, _ := tm.tokens.getAtSafe(0)
			_, _ = tm.identifyKeyword(tkn)
			assert.Equal(t, tt.expectedTknValue, tkn.getValue())
			t.Log(tm.tokens.sPrint())
		})
	}

}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func setUp() *tokens {
	return &tokens{
		newToken("A"),
		newToken("B"),
		newToken("C"),
		newToken("D"),
		newToken("E"),
	}
}

func TestTokenFunctions(t *testing.T) {
	// Test set up
	tokensList := setUp()

	// Test getAt function
	tkn, ok := tokensList.getAtSafe(1)
	assert.True(t, ok)
	assert.Equal(t, "B", tkn.getValue())

	// Reset tokens
	tokensList = setUp()

	// Test removeAt function
	tokensList.removeAt(1)
	tkn, ok = tokensList.getAtSafe(1)
	assert.True(t, ok)
	assert.Equal(t, "C", tkn.getValue())

	// Reset tokens
	tokensList = setUp()

	// Test overwriteAt function
	tokensList.overwriteAt(1, *newToken("X"))
	tkn, ok = tokensList.getAtSafe(1)
	assert.True(t, ok)
	assert.Equal(t, "X", tkn.getValue())

	// Reset tokens
	tokensList = setUp()

	// Test getFrom function
	tokensFrom := tokensList.getFrom(1)
	assert.Len(t, tokensFrom, 4)
	assert.Equal(t, "B", tokensFrom[0].getValue())
	assert.Equal(t, "C", tokensFrom[1].getValue())
	assert.Equal(t, "D", tokensFrom[2].getValue())
	assert.Equal(t, "E", tokensFrom[3].getValue())

	// Reset tokens
	tokensList = setUp()

	// Test getTo function
	tokensTo := tokensList.getTo(3)
	assert.Len(t, tokensTo, 3)
	assert.Equal(t, "A", tokensTo[0].getValue())
	assert.Equal(t, "B", tokensTo[1].getValue())
	assert.Equal(t, "C", tokensTo[2].getValue())

	// Reset tokens
	tokensList = setUp()

	// Test getToInc function
	tokensToInc := tokensList.getToInc(3)
	assert.Len(t, tokensToInc, 4)
	assert.Equal(t, "A", tokensToInc[0].getValue())
	assert.Equal(t, "B", tokensToInc[1].getValue())
	assert.Equal(t, "C", tokensToInc[2].getValue())
	assert.Equal(t, "D", tokensToInc[3].getValue())

	// Reset tokens
	tokensList = setUp()

	// Test getFromTo function
	tokensFromTo := tokensList.getFromTo(1, 3)
	assert.Len(t, tokensFromTo, 2)
	assert.Equal(t, "B", tokensFromTo[0].getValue())
	assert.Equal(t, "C", tokensFromTo[1].getValue())

	// Reset tokens
	tokensList = setUp()

	// Test getFromToInc function
	tokensFromToInc := tokensList.getFromToInc(0, 2)
	assert.Len(t, tokensFromToInc, 3)
	assert.Equal(t, "A", tokensFromToInc[0].getValue())
	assert.Equal(t, "B", tokensFromToInc[1].getValue())
	assert.Equal(t, "C", tokensFromToInc[2].getValue())

	// Test for insertAtStart method
	tokensList = setUp()
	tokensList.insertAtStart(*newToken("X"))
	tkn, _ = tokensList.getAtSafe(0)
	assert.Equal(t, "X", tkn.getValue())

	// Test for insertManyAt method
	tokensList = setUp()
	tokensList.insertManyAt(1, []*token{newToken("X"), newToken("Y")})
	tkn, _ = tokensList.getAtSafe(1)
	assert.Equal(t, "X", tkn.getValue())
	tkn, _ = tokensList.getAtSafe(2)
	assert.Equal(t, "Y", tkn.getValue())

	// Test for overwriteManyAt method
	tokensList = setUp()
	tokensList.overwriteManyAt(1, []*token{newToken("X"), newToken("Y")})
	tkn, _ = tokensList.getAtSafe(1)
	assert.Equal(t, "X", tkn.getValue())
	tkn, _ = tokensList.getAtSafe(2)
	assert.Equal(t, "Y", tkn.getValue())
}

func TestTokenSequenceFunctions(t *testing.T) {

	tm := newTokenManager("01 - 05")

	// Test getCategorySequenceAfter
	_, found := tm.tokens.getCategorySequenceAfter(0, []tokenCategory{
		tokenCatDelimiter, //
		tokenCatSeparator, // -
		tokenCatDelimiter, //
		tokenCatUnknown,   // 05
	}, false)
	assert.True(t, found)

	// Test getCategorySequenceAfterInc
	_, found = tm.tokens.getCategorySequenceAfterInc(0, []tokenCategory{
		tokenCatUnknown,   // 01
		tokenCatDelimiter, //
		tokenCatSeparator, // -
		tokenCatDelimiter, //
		tokenCatUnknown,   // 05
	}, false)
	assert.True(t, found)

	// Test getCategorySequenceAfter by skipping delimiters
	_, found = tm.tokens.getCategorySequenceAfter(0, []tokenCategory{
		tokenCatSeparator, // -
		tokenCatUnknown,   // 05
	}, true)
	assert.True(t, found)

	// Test getCategorySequenceAfterInc by skipping delimiters
	_, found = tm.tokens.getCategorySequenceAfterInc(0, []tokenCategory{
		tokenCatUnknown,   // 01
		tokenCatSeparator, // -
		tokenCatUnknown,   // 05
	}, true)
	assert.True(t, found)

	// Test getCategorySequenceBefore
	_, found = tm.tokens.getCategorySequenceBefore(4, []tokenCategory{
		tokenCatDelimiter, //
		tokenCatSeparator, // -
		tokenCatDelimiter, //
		tokenCatUnknown,   // 01
	}, false)
	assert.True(t, found)

	// Test getCategorySequenceBeforeInc
	_, found = tm.tokens.getCategorySequenceBeforeInc(4, []tokenCategory{
		tokenCatUnknown,   // 05
		tokenCatDelimiter, //
		tokenCatSeparator, // -
		tokenCatDelimiter, //
		tokenCatUnknown,   // 01
	}, false)
	assert.True(t, found)

	// Test getCategorySequenceBefore
	_, found = tm.tokens.getCategorySequenceBefore(4, []tokenCategory{
		tokenCatSeparator, // -
		tokenCatUnknown,   // 01
	}, true)
	assert.True(t, found)

	//

	_, found = tm.tokens.peekValuesAfter(0, []string{" ", "-", " ", "05"})
	assert.True(t, found)

	_, found = tm.tokens.peekValuesAfter(0, []string{" ", "-", " ", "05", " "}) // out of range
	assert.False(t, found)

}
