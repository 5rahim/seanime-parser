package seanime_parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Test set up
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

	// Test getCategorySequenceFrom
	_, found := tm.tokens.getCategorySequenceAfter(0, []tokenCategory{
		tokenCatDelimiter, //
		tokenCatSeparator, // -
		tokenCatDelimiter, //
		tokenCatUnknown,   // 05
	}, false)
	assert.True(t, found)

	// Test getCategorySequenceFrom by skipping delimiters
	_, found = tm.tokens.getCategorySequenceAfter(0, []tokenCategory{
		tokenCatSeparator, // -
		tokenCatUnknown,   // 05
	}, true)
	assert.True(t, found)

}
