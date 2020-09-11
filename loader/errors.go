package loader

import (
	"fmt"
	"strings"
)

type DefaultLanguageNotFoundError struct {
	Found   []string
	Default string
}

func (err *DefaultLanguageNotFoundError) Error() string {
	return fmt.Sprintf("default language '%s' not found in %s",
		err.Default,
		strings.Join(err.Found, ","),
	)
}

type KeysNotFoundInDefaultError struct {
	Language        string
	DefaultLanguage string
	Keys            []string
}

func (err *KeysNotFoundInDefaultError) Error() string {
	return fmt.Sprintf("langugage '%s' has key(s) (%s) not present in default language '%s'",
		err.Language,
		strings.Join(err.Keys, ","),
		err.DefaultLanguage,
	)
}
