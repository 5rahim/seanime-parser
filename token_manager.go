package seanime_parser

import (
	"fmt"
	"strings"
)

type tokens []*token

type tokenManager struct {
	tokens         *tokens
	keywordManager *keywordManager
	filename       string
}

func newTokenManager(filename string) *tokenManager {
	tm := tokenManager{
		tokens:         &tokens{},
		filename:       filename,
		keywordManager: newKeywordManager(),
	}

	tm.tokens.setTokens(tokenize(strings.TrimSpace(filename)))

	tm.mergeDecimals()

	return &tm
}

func (tm *tokenManager) mergeDecimals() {
	for _, tkn := range *tm.tokens {
		if !tkn.isNumberKind() {
			continue
		}

		_, _ = tm.tokens.checkNumberWithDecimal(tkn)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// checkNumberWithDecimal checks if a token (number) is followed by a decimal point and a number.
// If it is, it will merge the tokens into a single token and return it.
func (t *tokens) checkNumberWithDecimal(tkn *token) (*token, bool) {
	if tkn == nil || !tkn.isNumberKind() {
		return nil, false
	}

	// Check if token is followed by a decimal point and a number
	dotTkn, ok := t.getTokenAfter(tkn)
	if !ok || !dotTkn.isDotDelimiter() {
		return nil, false
	}

	numTkn, ok := t.getTokenAfter(dotTkn)
	if !ok || !numTkn.isNumberKind() || (numTkn.getValue() != "5" && numTkn.getValue() != "05") {
		return nil, false
	}

	delTkn, ok := t.getTokenAfter(numTkn)
	if (!ok || !delTkn.isDelimiter()) && !t.isLastToken(numTkn) { // Delimiter or end of tokens
		return nil, false
	}

	// Merge tokens
	tkn.setValue(tkn.getValue() + "." + numTkn.getValue())

	// Remove dot and number tokens
	t.removeAt(t.getIndexOf(dotTkn))
	t.removeAt(t.getIndexOf(numTkn))

	return tkn, true
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// getTokenAfter returns the token that comes after the specified token, along with a boolean indicating if the token was found.
func (t *tokens) getTokenAfter(tkn *token) (*token, bool) {
	index := t.getIndexOf(tkn)
	if index == -1 {
		return nil, false
	}
	return t.getAtSafe(index + 1)
}

// getTokenAfterSD returns the token that comes after the specified token, along with a boolean indicating if the token was found.
// It searches for the next non-delimiter token in the tokens slice starting from the index of the specified token.
// It also returns the number of skipped delimiter tokens.
func (t *tokens) getTokenAfterSD(tkn *token) (*token, bool, int) {
	index := t.getIndexOf(tkn)
	if index == -1 {
		return nil, false, 0
	}
	skipped := 0
	for i := index + 1; i < len(*t); i++ {
		if !(*t)[i].isDelimiter() {
			return (*t)[i], true, skipped
		} else {
			skipped++
		}
	}
	return nil, false, skipped
}

// getTokenBefore returns the token that comes before the specified token, along with a boolean indicating if the token was found.
func (t *tokens) getTokenBefore(tkn *token) (*token, bool) {
	index := t.getIndexOf(tkn)
	if index == -1 {
		return nil, false
	}
	return t.getAtSafe(index - 1)
}

// getTokenBeforeSD returns the token that comes before the specified token, along with a boolean indicating if the token was found.
// It searches for the previous non-delimiter token in the tokens slice starting from the index of the specified token.
// It also returns the number of skipped delimiter tokens.
func (t *tokens) getTokenBeforeSD(tkn *token) (*token, bool, int) {
	index := t.getIndexOf(tkn)
	if index == -1 {
		return nil, false, 0
	}
	skipped := 0
	for i := index - 1; i >= 0; i-- {
		if !(*t)[i].isDelimiter() {
			return (*t)[i], true, skipped
		} else {
			skipped++
		}
	}
	return nil, false, skipped
}

// isTokenInFirstHalf checks if the specified token is in the first half of the tokens list.
// It returns true if the token is found and its index is less than or equal to half the length of the list,
// otherwise it returns false.
func (t *tokens) isTokenInFirstHalf(tkn *token) bool {
	index := t.getIndexOf(tkn)
	if index == -1 {
		return false
	}
	return index <= len(*t)/2
}

// isTokenAfterFileMetadata checks if the specified token comes after file info metadata
func (t *tokens) isTokenAfterFileMetadata(tkn *token) bool {
	index := t.getIndexOf(tkn)
	if index == -1 {
		return false
	}
	isAfter := false

	for idx, _tkn := range *t {
		// Check if token is after file info metadata
		if _tkn.isFileInfoMetadata() && idx != index && idx < index {
			isAfter = true
		}
	}

	return isAfter
}

func (t *tokens) getIndexOf(tkn *token) int {
	for i, _tkn := range *t {
		if _tkn.UUID == tkn.UUID {
			return i
		}
	}
	return -1
}

func (t *tokens) isLastToken(tkn *token) bool {
	if tkn == nil {
		return false
	}
	return t.getIndexOf(tkn) == len(*t)-1
}

func (t *tokens) foundDashSeparatorBefore(tkn *token) bool {
	// Check if token before previous token is a dash separator
	if prevPrevTkn, found, _ := t.getTokenBeforeSD(tkn); found {
		if !prevPrevTkn.isDashSeparator() {
			return false
		}
	} else {
		return false
	}
	return true
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (t *tokens) setTokens(tkns []*token) {
	*t = tkns
}

func (t *tokens) insertAt(index int, tkn token) {
	if index < 0 || index > len(*t) {
		return
	}
	*t = append((*t)[:index], append([]*token{&tkn}, (*t)[index:]...)...)
}

func (t *tokens) insertAtEnd(tkn token) {
	*t = append(*t, &tkn)
}

func (t *tokens) insertAtStart(tkn token) {
	*t = append([]*token{&tkn}, *t...)
}

func (t *tokens) insertManyAt(index int, tkns []*token) {
	if index < 0 || index > len(*t) {
		return
	}
	*t = append((*t)[:index], append(tkns, (*t)[index:]...)...)
}

func (t *tokens) insertAfter(index int, tkn token) {
	if index < 0 || index > len(*t) {
		return
	}
	*t = append((*t)[:index+1], append([]*token{&tkn}, (*t)[index+1:]...)...)
}

func (t *tokens) insertManyAfter(index int, tkns []*token) {
	if index < 0 || index > len(*t) {
		return
	}
	*t = append((*t)[:index+1], append(tkns, (*t)[index+1:]...)...)
}

func (t *tokens) removeAt(index int) {
	if index < 0 || index > len(*t) {
		return
	}
	*t = append((*t)[:index], (*t)[index+1:]...)
}

func (t *tokens) overwriteAt(index int, tkn token) {
	(*t)[index] = &tkn
}

func (t *tokens) overwriteManyAt(index int, tkns []*token) {
	*t = append((*t)[:index], append(tkns, (*t)[index+len(tkns):]...)...)
}

func (t *tokens) overwriteAndInsertManyAt(index int, tkns []*token) {
	*t = append((*t)[:index], (*t)[index+1:]...)
	// Then insert new elements at index
	// append takes a slice and follows that with a variadic parameter hence the need for ...
	*t = append((*t)[:index], append(tkns, (*t)[index:]...)...)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (t *tokens) getAtSafe(index int) (*token, bool) {
	if index < 0 || index > len(*t)-1 {
		return nil, false
	}
	return (*t)[index], true
}
func (t *tokens) getAt(index int) *token {
	return (*t)[index]
}

func (t *tokens) getFromUUID(uuid string) *token {
	for _, tkn := range *t {
		if tkn.UUID == uuid {
			return tkn
		}
	}
	return nil
}

func (t *tokens) getFromUUIDSafe(uuid string) (*token, bool) {
	for _, tkn := range *t {
		if tkn.UUID == uuid {
			return tkn, true
		}
	}
	return nil, false
}

func (t *tokens) getFromUUIDs(uuids []string) []*token {
	tkns := make([]*token, 0)
	for _, uuid := range uuids {
		tkn := t.getFromUUID(uuid)
		if tkn != nil {
			tkns = append(tkns, tkn)
		}
	}
	return tkns
}

func (t *tokens) getFrom(index int) []*token {
	if index < 0 || index > len(*t) {
		return []*token{}
	}
	return (*t)[index:]
}

func (t *tokens) getTo(index int) []*token {
	if index < 0 || index > len(*t) {
		return []*token{}
	}
	return (*t)[:index]
}

func (t *tokens) getToInc(index int) []*token {
	if index < 0 || index+1 > len(*t) {
		return []*token{}
	}
	return (*t)[:index+1]
}

func (t *tokens) getFromTo(start int, end int) []*token {
	// check indices
	if start < 0 || end < 0 || start > end || start > len(*t) || end > len(*t) {
		return []*token{}
	}
	return (*t)[start:end]
}

func (t *tokens) getFromToInc(start int, end int) []*token {
	// check indices
	if start < 0 || end < 0 || start > end || start+1 > len(*t) || end+1 > len(*t) {
		return []*token{}
	}
	return (*t)[start : end+1]
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (t *tokens) getFirstOccurrenceAfter(start int, pred func(tkn *token) bool) (*token, bool) {
	if start < 0 {
		start = -1
	}
	if start+1 > len(*t) {
		return nil, false
	}
	for i := start + 1; i < len(*t); i++ {
		if pred((*t)[i]) {
			return (*t)[i], true
		}
	}
	return nil, false
}

func (t *tokens) getFirstOccurrenceBefore(start int, pred func(tkn *token) bool) (*token, bool) {
	if start > len(*t) {
		start = len(*t) + 1
	}
	if start < 0 {
		return nil, false
	}
	for i := start - 1; i >= 0; i-- {
		if pred((*t)[i]) {
			return (*t)[i], true
		}
	}
	return nil, false
}

// getCategorySequenceAfter returns the sequence of tokens in the given categories after the specified start index,
// along with a boolean indicating if the sequence was found.
// The skipDelimiters parameter determines whether to skip delimiter tokens when collecting the sequence.
func (t *tokens) getCategorySequenceAfter(start int, categories []tokenCategory, skipDelimiters bool) ([]*token, bool, int) {
	if start < 0 {
		start = -1
	}
	if start+1 > len(*t) {
		return []*token{}, false, 0
	}

	nbSkipped := 0

	var collec []*token
	var cursor int
	for i := start + 1; i < len(*t); i++ {
		if len(collec) == len(categories) {
			break
		}
		if skipDelimiters && (*t)[i].isDelimiter() {
			nbSkipped += 1
			continue
		}
		if (*t)[i].isCategory(categories[cursor]) {
			collec = append(collec, (*t)[i])
			cursor++
		} else {
			break
		}
	}

	if len(collec) == len(categories) {
		return collec, true, nbSkipped
	}

	return []*token{}, false, 0
}

func (t *tokens) getCategorySequenceAfterInc(start int, categories []tokenCategory, skipDelimiters bool) ([]*token, bool, int) {
	return t.getCategorySequenceAfter(start-1, categories, skipDelimiters)
}

// getCategorySequenceBefore returns the sequence of tokens in the given categories before the specified start index,
// along with a boolean indicating if the sequence was found.
// The skipDelimiters parameter determines whether to skip delimiter tokens when collecting the sequence.
func (t *tokens) getCategorySequenceBefore(start int, categories []tokenCategory, skipDelimiters bool) ([]*token, bool, int) {
	if start > len(*t) {
		start = len(*t) + 1
	}
	if start < 0 {
		return []*token{}, false, 0
	}

	nbSkipped := 0

	var collec []*token
	var cursor int
	for i := start - 1; i >= 0; i-- {
		if len(collec) == len(categories) {
			break
		}
		if skipDelimiters && (*t)[i].isDelimiter() {
			nbSkipped += 1
			continue
		}
		if (*t)[i].isCategory(categories[cursor]) {
			collec = append([]*token{(*t)[i]}, collec...)
			cursor++
		} else {
			break
		}
	}

	if len(collec) == len(categories) {
		return collec, true, nbSkipped
	}

	return []*token{}, false, 0
}

func (t *tokens) getCategorySequenceBeforeInc(start int, categories []tokenCategory, skipDelimiters bool) ([]*token, bool, int) {
	return t.getCategorySequenceBefore(start+1, categories, skipDelimiters)
}

func (t *tokens) iterate(iterationFunc func(tkn *token, idx int)) {
	for idx, tkn := range *t {
		iterationFunc(tkn, idx)
	}
}

////////////////////

func (t *tokens) peekValuesAfter(start int, strs []string) ([]*token, bool) {

	if start+1+len(strs) > len(*t) {
		return nil, false
	}

	_tkns := (*t)[start+1 : start+1+len(strs)]

	var collec []*token
	for i := 0; i < len(strs); i++ {
		if strings.ToUpper(_tkns[i].getValue()) == strings.ToUpper(strs[i]) {
			uuid, ok := t.getFromUUIDSafe(_tkns[i].UUID)
			if !ok {
				break
			}
			collec = append(collec, uuid)
		} else {
			break
		}
	}

	if len(collec) == len(strs) {
		return collec, true
	}

	return nil, false
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (t *tokens) findWithMetadataKind(cat metadataCategory) (bool, []*token) {
	_tkns := make([]*token, 0)
	for _, tkn := range *t {
		if tkn.MetadataCategory == cat {
			_tkns = append(_tkns, tkn)
		}
	}
	if len(_tkns) > 0 {
		return true, _tkns
	}
	return false, nil
}

func (t *tokens) findWithMetadataCategory(cat tokenCategory) (bool, []*token) {
	_tkns := make([]*token, 0)
	for _, tkn := range *t {
		if tkn.isCategory(cat) {
			_tkns = append(_tkns, tkn)
		}
	}
	if len(_tkns) > 0 {
		return true, _tkns
	}
	return false, nil
}

func (t *tokens) findWithKeywordCategory(cat keywordCategory) (bool, []*token) {
	_tkns := make([]*token, 0)
	for _, tkn := range *t {
		if tkn.IdentifiedKeywordCategory == cat && tkn.isKeyword() {
			_tkns = append(_tkns, tkn)
		}
	}
	if len(_tkns) > 0 {
		return true, _tkns
	}
	return false, nil
}

func (t *tokens) sPrint() string {
	str := "["
	for idx, tkn := range *t {
		str += "\"" + tkn.getValue()
		if idx < len(*t)-1 {
			str += "\", "
		} else {
			str += "\""
		}
	}
	str += "]"
	return str
}

func (t *tokens) sDump() string {
	str := "\n"
	for _, tkn := range *t {
		str += fmt.Sprintf("%-12s\t%v, kw: %v, %v, m: %v\n",
			"\""+tkn.getValue()+"\"",
			tkn.getCategory(),
			tkn.IdentifiedKeywordCategory,
			tkn.getKind(),
			tkn.MetadataCategory,
		)
	}
	str += "\n"
	return str
}
