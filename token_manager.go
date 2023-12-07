package seanime_parser

import (
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

	return &tm
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (t *tokens) getTokenAfter(tkn *token) (*token, bool) {
	index := t.getIndexOf(tkn)
	if index == -1 {
		return nil, false
	}
	return t.getAtSafe(index + 1)
}

func (t *tokens) getTokenBefore(tkn *token) (*token, bool) {
	index := t.getIndexOf(tkn)
	if index == -1 {
		return nil, false
	}
	return t.getAtSafe(index - 1)
}

func (t *tokens) isTokenInFirstHalf(tkn *token) bool {
	index := t.getIndexOf(tkn)
	if index == -1 {
		return false
	}
	return index <= len(*t)/2
}

func (t *tokens) isTokenAfterFileMetadata(tkn *token) bool {
	index := t.getIndexOf(tkn)
	if index == -1 {
		return false
	}
	isAfter := false

	for idx, _tkn := range *t {
		// Check if token is after file info metadata
		if _tkn.isFileInfoMetadata() && idx != index && idx > index {
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

func (t *tokens) getAtSafe(index int) (*token, bool) {
	if index < 0 || index > len(*t) {
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

func (t *tokens) getFirstCategoryOccurrence(cat tokenCategory) (*token, bool) {
	for _, tkn := range *t {
		if tkn.Category == cat {
			return tkn, true
		}
	}
	return nil, false
}

func (t *tokens) getLastCategoryOccurrence(cat tokenCategory) (*token, bool) {
	for i := len(*t) - 1; i >= 0; i-- {
		if (*t)[i].Category == cat {
			return (*t)[i], true
		}
	}
	return nil, false
}

func (t *tokens) getFirstCategoryOccurrenceAfter(start int, cat tokenCategory) (*token, bool) {
	if start < 0 {
		start = -1
	}
	if start+1 > len(*t) {
		return nil, false
	}
	for i := start + 1; i < len(*t); i++ {
		if (*t)[i].Category == cat {
			return (*t)[i], true
		}
	}
	return nil, false
}

func (t *tokens) getFirstCategoryOccurrenceBefore(start int, cat tokenCategory) (*token, bool) {
	if start > len(*t) {
		start = len(*t) + 1
	}
	if start < 0 {
		return nil, false
	}
	for i := start - 1; i >= 0; i-- {
		if (*t)[i].Category == cat {
			return (*t)[i], true
		}
	}
	return nil, false
}

// getCategorySequenceAfter returns the sequence of tokens in the given categories after the specified start index,
// along with a boolean indicating if the sequence was found.
// The skipDelimiters parameter determines whether to skip delimiter tokens when collecting the sequence.
func (t *tokens) getCategorySequenceAfter(start int, categories []tokenCategory, skipDelimiters bool) ([]*token, bool) {
	if start < 0 {
		start = -1
	}
	if start+1 > len(*t) {
		return []*token{}, false
	}

	_tkns := make([]*token, 0)
	if skipDelimiters {
		for i := start + 1; i < len(*t); i++ {
			if !(*t)[i].isDelimiter() {
				_tkns = append(_tkns, (*t)[i])
			} else {
				continue
			}
		}
	} else {
		_tkns = (*t)[start+1:]
	}

	var collec []*token
	for i := 0; i < len(categories); i++ {
		if _tkns[i].isCategory(categories[i]) {
			uuid, ok := t.getFromUUIDSafe(_tkns[i].UUID)
			if !ok {
				break
			}
			collec = append(collec, uuid)
		} else {
			break
		}
	}

	if len(collec) == len(categories) {
		return collec, true
	}

	return []*token{}, false
}

func (t *tokens) getCategorySequenceAfterInc(start int, categories []tokenCategory, skipDelimiters bool) ([]*token, bool) {
	return t.getCategorySequenceAfter(start-1, categories, skipDelimiters)
}

// getCategorySequenceBefore returns the sequence of tokens in the given categories before the specified start index,
// along with a boolean indicating if the sequence was found.
// The skipDelimiters parameter determines whether to skip delimiter tokens when collecting the sequence.
func (t *tokens) getCategorySequenceBefore(start int, categories []tokenCategory, skipDelimiters bool) ([]*token, bool) {
	if start > len(*t) {
		start = len(*t) + 1
	}
	if start < 0 {
		return []*token{}, false
	}

	_tkns := make([]*token, 0)
	if skipDelimiters {
		for i := start - 1; i >= 0; i-- {
			if !(*t)[i].isDelimiter() {
				_tkns = append(_tkns, (*t)[i])
			} else {
				continue
			}
		}
	} else {
		for i := start - 1; i >= 0; i-- {
			_tkns = append(_tkns, (*t)[i])
		}
	}

	var collec []*token
	for i := 0; i < len(categories); i++ {
		if _tkns[i].isCategory(categories[i]) {
			uuid, ok := t.getFromUUIDSafe(_tkns[i].UUID)
			if !ok {
				break
			}
			collec = append(collec, uuid)
		} else {
			break
		}
	}

	if len(collec) == len(categories) {
		return collec, true
	}

	return []*token{}, false
}

func (t *tokens) getCategorySequenceBeforeInc(start int, categories []tokenCategory, skipDelimiters bool) ([]*token, bool) {
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
