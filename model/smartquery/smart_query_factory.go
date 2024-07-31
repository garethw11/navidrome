package smartquery

import (
	"errors"
	"strings"

	"github.com/navidrome/navidrome/log"
)

type SmartQueryFactory struct {
	Definition []string
}

func (r *SmartQueryFactory) CreateSmartQuery() (*SmartQuery, error) {
	// first line must be PLAYLIST name:your_name_here, description:describe_your_query
	firstLine := strings.TrimSpace(r.Definition[0])
	if !strings.HasPrefix(firstLine, "PLAYLIST ") {
		return nil, errors.New("smart query first line must be \"PLAYLIST name:your_name_here, description:describe_your_query\" but " +
			"you supplied \"" + strings.TrimSpace(r.Definition[0]) + "\"")
	}
	firstLine = strings.TrimPrefix(firstLine, "PLAYLIST ")
	name, description, success := r.extractNameAndDescription(firstLine)
	if !success {
		return nil, errors.New("smart query first line must be \"PLAYLIST name:your_name_here, description:describe_your_query\" but " +
			"you supplied \"" + strings.TrimSpace(r.Definition[0]) + "\"")
	}
	r.Definition = r.Definition[1:]
	query := strings.Join(r.Definition, "\n")
	log.Debug("Name %v Description %v", name, description)
	log.Debug("Query [" + query + "]")
	return &SmartQuery{Name: name, Comment: description, Query: query, OrderBy: "title"}, nil
}

// TODO this doesn't support , in name or description.  Do we want commas or use different char eg | or just ban commas?
func (r *SmartQueryFactory) extractNameAndDescription(firstLine string) (string, string, bool) {
	keyValues := strings.Split(strings.TrimSpace(firstLine), ",")
	if len(keyValues) != 2 {
		return "", "", false
	}
	name, isNameFound := r.extractValue(keyValues[0], "name:")
	description, isDescriptionFound := r.extractValue(keyValues[1], "description:")
	return name, description, isNameFound || isDescriptionFound
}

func (r *SmartQueryFactory) extractValue(text string, key string) (string, bool) {
	text = strings.TrimSpace(text)
	if !strings.HasPrefix(text, key) {
		return "", false
	}
	value := strings.TrimSpace(strings.TrimPrefix(text, key))
	if value == "" {
		return "", false
	}
	return value, true
}
