package executor

import (
	"regexp"

	"github-release-puller/config"
)

type filters []*regexp.Regexp

func compile(raw []config.AssetFilter) (filters, error) {
	result := make([]*regexp.Regexp, 0, len(raw))

	for _, r := range raw {
		regex, err := regexp.Compile(r.Match)
		if err != nil {
			return nil, err
		}

		result = append(result, regex)
	}

	return result, nil
}

func (f filters) match(label string) bool {
	for _, filter := range f {
		if filter.MatchString(label) {
			return true
		}
	}

	return false
}
