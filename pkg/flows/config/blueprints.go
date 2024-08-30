package config

import (
	_ "embed"
	"fmt"
)

//go:embed blueprints/default.yaml
var BlueprintDefault []byte

//go:embed blueprints/context.yaml
var BlueprintContext []byte

var Blueprints = map[string][]byte{
	"default": BlueprintDefault,
	"context": BlueprintContext,
}

func GetBlueprint(name string) ([]byte, error) {
	bp, ok := Blueprints[name]
	if !ok {
		return nil, fmt.Errorf("blueprint %q not found", name)
	}
	return bp, nil
}
