package seanime_parser

import (
	"github.com/google/uuid"
	"strings"
)

type tokens []*token

type tokenManager struct {
	tokens         *tokens
	keywordManager *keywordManager
	filename       string
}

func newTokenManager(filename string) tokenManager {
	tm := tokenManager{
		tokens:         &tokens{},
		filename:       filename,
		keywordManager: newKeywordManager(),
	}

	tm.tokens.setTokens(tokenize(strings.TrimSpace(filename)))

	return tm
}

func (t *tokens) setTokens(tkns []*token) {
	*t = tkns
}

func (t *tokens) insert(index int, tkn token) {
	if index == 0 {
		if len(*t) == 0 {
			tkn.UUID = uuid.New().String()
			*t = append(*t, &tkn)
			return
		} else if len(*t) == 1 && tkn.Value != (*t)[index].Value {
			tkn.UUID = uuid.New().String()
			(*t)[index] = &tkn
			return
		}
	} else if index > len(*t)-1 {
		return
	} else if index < 0 {
		return
	}
	if (*t)[index].Value == tkn.Value {
		return
	}
	tkn.UUID = uuid.New().String()
	startList := append((*t)[:index], &tkn)
	*t = append(startList, (*t)[index:]...)
}

func (tm *tokenManager) getTokens() []*token {
	return *tm.tokens
}
