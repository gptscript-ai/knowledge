package config

import (
	"encoding/json"
	"github.com/gptscript-ai/knowledge/pkg/datastore/documentloader"
	"github.com/gptscript-ai/knowledge/pkg/flows"
	"os"
	"sigs.k8s.io/yaml"
	"strings"
)

type FlowConfig struct {
	Flows map[string]FlowConfigEntry `json:"flows" yaml:"flows" mapstructure:"flows"`
}

type DocumentLoaderConfig struct {
	Name    string         `json:"name" yaml:"name" mapstructure:"name"`
	Options map[string]any `json:"options,omitempty" yaml:"options" mapstructure:"options"`
}

type TextSplitterConfig struct {
	Name    string         `json:"name" yaml:"name" mapstructure:"name"`
	Options map[string]any `json:"options,omitempty" yaml:"options" mapstructure:"options"`
}

type IngestionFlowConfig struct {
	Filetypes      []string             `json:"filetypes" yaml:"filetypes" mapstructure:"filetypes"`
	DocumentLoader DocumentLoaderConfig `json:"documentLoader,omitempty" yaml:"documentLoader" mapstructure:"documentLoader"`
	TextSplitter   TextSplitterConfig   `json:"textSplitter,omitempty" yaml:"textSplitter" mapstructure:"textSplitter"`
	Transformers   []string             `json:"transformers,omitempty" yaml:"transformers" mapstructure:"transformers"`
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

func (i *IngestionFlowConfig) AsIngestionFlow() (*flows.IngestionFlow, error) {
	flow := &flows.IngestionFlow{}
	if i.DocumentLoader.Name != "" {
		name := strings.ToLower(strings.Trim(i.DocumentLoader.Name, " "))
		cfg, err := documentloader.GetDocumentLoaderConfig(name)
		if err != nil {
			return nil, err
		}
		if len(i.DocumentLoader.Options) > 0 {
			jsondata, err := json.Marshal(i.DocumentLoader.Options)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(jsondata, &cfg)
			if err != nil {
				return nil, err
			}
		}
	}

	return flow, nil
}
