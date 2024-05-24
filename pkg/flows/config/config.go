package config

import (
	"encoding/json"
	"os"
	"sigs.k8s.io/yaml"
)

type FlowConfig struct {
	Flows map[string]FlowConfigEntry `json:"flows" yaml:"flows" mapstructure:"flows"`
}

type IngestionFlowConfig struct {
	Filetypes      []string `json:"filetypes" yaml:"filetypes" mapstructure:"filetypes"`
	DocumentLoader string   `json:"documentLoader,omitempty" yaml:"documentLoader" mapstructure:"documentLoader"`
	TextSplitter   string   `json:"textSplitter,omitempty" yaml:"textSplitter" mapstructure:"textSplitter"`
	Transformers   []string `json:"transformers,omitempty" yaml:"transformers" mapstructure:"transformers"`
}

type RetrievalFlowConfig struct{}

type FlowConfigEntry struct {
	Ingestion []IngestionFlowConfig `json:"ingestion,omitempty" yaml:"ingestion" mapstructure:"ingestion"`
	Retrieval RetrievalFlowConfig   `json:"retrieval,omitempty" yaml:"retrieval" mapstructure:"retrieval"`
}

// FromFile reads a configuration file and returns a FlowConfig.
func FromFile(filename string) (*FlowConfig, error) {

	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config FlowConfig
	jsondata := content
	if !json.Valid(content) {
		jsondata, err = yaml.YAMLToJSON(content)
		if err != nil {
			return nil, err
		}
	}

	err = json.Unmarshal(jsondata, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
