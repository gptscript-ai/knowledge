package load

import (
	"fmt"
	"path"
	"strings"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	koanf "github.com/knadh/koanf/v2"
)

var k = koanf.New(".")

// FillConfig fills the given struct with values from a config file and environment variables.
// The envPrefix parameter is used to prefix environment variables.
func FillConfig(configFile string, envPrefix string, cfg interface{}) error {

	if configFile != "" {
		// yaml or json
		var pa koanf.Parser
		switch path.Ext(configFile) {
		case ".json":
			pa = json.Parser()
		case ".yaml", ".yml":
			pa = yaml.Parser()
		default:
			return fmt.Errorf("unsupported config file format: %s", path.Ext(configFile))
		}

		if err := k.Load(file.Provider(configFile), pa); err != nil {
			return fmt.Errorf("error loading config file: %w", err)
		}
	}

	// Load environment variables and override config file values
	envProvider := env.Provider(envPrefix, ".", func(s string) string {
		return strings.Replace(strings.ToLower(s), "_", "-", -1)
	})
	if err := k.Load(envProvider, nil); err != nil {
		return fmt.Errorf("error loading environment variables: %w", err)
	}

	// Unmarshal the configuration into the provided struct
	if err := k.Unmarshal("", cfg); err != nil {
		return fmt.Errorf("error unmarshaling config: %w", err)
	}

	return nil
}
