package loader

import (
	"regexp"
)

type stringSet map[string]struct{}

func newStringSet(in []string) stringSet {
	result := make(stringSet, len(in))
	for _, v := range in {
		result[v] = struct{}{}
	}
	return result
}

func (ss stringSet) Diff(ss2 stringSet) stringSet {
	result := make(stringSet, len(ss))

	for k := range ss {
		if _, exists := ss2[k]; !exists {
			result[k] = struct{}{}
		}
	}

	return result
}

func (ss stringSet) Slice() []string {
	result := make([]string, len(ss))
	i := 0
	for k := range ss {
		result[i] = k
		i++
	}
	return result
}

var templateVarsRegexp = regexp.MustCompile(`\{(\w+)\}`)

func getTemplateVars(v string) []string {
	matches := templateVarsRegexp.FindAllStringSubmatch(v, -1)

	args := make([]string, len(matches))

	for i, match := range matches {
		args[i] = match[1]
	}

	return args
}
