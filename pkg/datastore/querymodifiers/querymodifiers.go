package querymodifiers

import (
	"fmt"
)

type QueryModifier interface {
	ModifyQuery(query string) (string, error)
}

var QueryModifiers = map[string]QueryModifier{
	"spellcheck": SpellcheckQueryModifier{},
}

func GetQueryModifier(name string) (QueryModifier, error) {
	qm, ok := QueryModifiers[name]
	if !ok {
		return nil, fmt.Errorf("unknown queryModifier %q", name)
	}
	return qm, nil
}
