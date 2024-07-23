package load

import (
	"fmt"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	koanf "github.com/knadh/koanf/v2"
	"path"
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
			return fmt.Errorf("error loading config file %q: %w", configFile, err)
		}
	}

	// Load environment variables and override config file values
	if err := k.Load(env.Provider(envPrefix, ".", func(s string) string { return s }), nil); err != nil {
		return fmt.Errorf("error loading environment variables: %w", err)
	}

	x := k.All()
	fmt.Printf("Config: %#v\n", x)

	// Unmarshal file config into the struct
	if err := k.UnmarshalWithConf("", cfg, koanf.UnmarshalConf{Tag: "koanf"}); err != nil {
		return fmt.Errorf("error unmarshalling file config: %w", err)
	}

	// Unmarshal environment variables into the struct, thus overriding file config values if present
	if err := k.UnmarshalWithConf("", cfg, koanf.UnmarshalConf{Tag: "env", FlatPaths: true}); err != nil {
		return fmt.Errorf("error unmarshalling environment variables: %w", err)
	}

	return nil
}
