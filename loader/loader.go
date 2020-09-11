package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Load all files matched glob from dir
func Load(dir string, glob string, defaultLang string) (*Bundle, error) {
	languages, err := loadLanguages(dir, glob)
	if err != nil {
		return nil, fmt.Errorf("load languages: %w", err)
	}

	bundle := &Bundle{
		All: languages,
	}

	// find default language
	bundle.Default = bundle.Get(defaultLang)
	if bundle.Default == nil {
		return nil, &DefaultLanguageNotFoundError{
			Default: defaultLang,
			Found:   bundle.Languages(),
		}
	}

	if err := validateKeysPresent(bundle); err != nil {
		return nil, fmt.Errorf("validate keys present: %w", err)
	}

	return bundle, nil
}

func validateKeysPresent(bundle *Bundle) error {
	keys := newStringSet(bundle.Default.Keys())

	for _, lang := range bundle.All {
		if lang.Name == bundle.Default.Name {
			continue
		}

		langKeys := newStringSet(lang.Keys())

		if diff := langKeys.Diff(keys); len(diff) > 0 {
			return &KeysNotFoundInDefaultError{
				DefaultLanguage: bundle.Default.Name,
				Language:        lang.Name,
				Keys:            diff.Slice(),
			}
		}
	}

	return nil
}

func loadLanguages(dir, glob string) ([]*Language, error) {
	path := filepath.Join(dir, glob)
	matches, err := filepath.Glob(path)
	if err != nil {
		return nil, fmt.Errorf("glob search: %w", err)
	}

	languages := make([]*Language, len(matches))
	i := 0
	for _, match := range matches {
		filename := filepath.Base(match)
		language := strings.TrimSuffix(filename, filepath.Ext(filename))

		messages, err := loadMessages(match)
		if err != nil {
			return nil, fmt.Errorf("load messages: %w", err)
		}

		languages[i] = &Language{
			Name:     language,
			File:     filename,
			Messages: messages,
		}
		i++
	}

	return languages, nil
}

func loadMessages(path string) ([]*Message, error) {
	dst := map[string]string{}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open lang file '%s': %w", path, err)
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(&dst); err != nil {
		return nil, fmt.Errorf("decode lang file '%s': %w", path, err)
	}

	result := make([]*Message, len(dst))

	{
		i := 0
		for k, v := range dst {
			result[i] = newMessage(k, v)
			i++
		}
	}

	return result, nil
}
