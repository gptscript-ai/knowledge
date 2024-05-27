package config

import (
	"encoding/json"
	"fmt"
	"github.com/gptscript-ai/knowledge/pkg/datastore/documentloader"
	"github.com/gptscript-ai/knowledge/pkg/datastore/textsplitter"
	"github.com/gptscript-ai/knowledge/pkg/flows"
	"os"
	"sigs.k8s.io/yaml"
	"strings"
)

type FlowConfig struct {
	Flows    map[string]FlowConfigEntry `json:"flows" yaml:"flows" mapstructure:"flows"`
	Datasets map[string]string          `json:"datasets,omitempty" yaml:"datasets" mapstructure:"datasets"`
}

type FlowConfigEntry struct {
	Default   bool                  `json:"default,omitempty" yaml:"default" mapstructure:"default"`
	Ingestion []IngestionFlowConfig `json:"ingestion,omitempty" yaml:"ingestion" mapstructure:"ingestion"`
	Retrieval *RetrievalFlowConfig  `json:"retrieval,omitempty" yaml:"retrieval" mapstructure:"retrieval"`
}

type IngestionFlowConfig struct {
	Filetypes      []string             `json:"filetypes" yaml:"filetypes" mapstructure:"filetypes"`
	DocumentLoader DocumentLoaderConfig `json:"documentLoader,omitempty" yaml:"documentLoader" mapstructure:"documentLoader"`
	TextSplitter   TextSplitterConfig   `json:"textSplitter,omitempty" yaml:"textSplitter" mapstructure:"textSplitter"`
	Transformers   []string             `json:"transformers,omitempty" yaml:"transformers" mapstructure:"transformers"`
}

type RetrievalFlowConfig struct{}

type DocumentLoaderConfig struct {
	Name    string         `json:"name" yaml:"name" mapstructure:"name"`
	Options map[string]any `json:"options,omitempty" yaml:"options" mapstructure:"options"`
}

type TextSplitterConfig struct {
	Name    string         `json:"name" yaml:"name" mapstructure:"name"`
	Options map[string]any `json:"options,omitempty" yaml:"options" mapstructure:"options"`
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

	return &config, config.Validate()
}

func (f *FlowConfig) Validate() error {
	defaultCount := 0
	for name, flow := range f.Flows {
		if flow.Default {
			defaultCount++
		}

		if len(flow.Ingestion) == 0 && flow.Retrieval == nil {
			return fmt.Errorf("flow %q has neither ingestion nor retrieval specified", name)
		}

	}
	if defaultCount > 1 {
		return fmt.Errorf("only one flow can be default, found %d", defaultCount)
	}
	return nil
}

func (f *FlowConfig) GetDefaultFlowConfigEntry() (*FlowConfigEntry, error) {
	for _, flow := range f.Flows {
		if flow.Default {
			return &flow, nil
		}
	}
	return nil, fmt.Errorf("default flow not found")
}

func (f *FlowConfig) GetFlow(name string) (*FlowConfigEntry, error) {
	flow, ok := f.Flows[name]
	if !ok {
		return nil, fmt.Errorf("flow %q not found", name)
	}
	return &flow, nil
}

// AsIngestionFlow converts an IngestionFlowConfig to an actual flows.IngestionFlow.
func (i *IngestionFlowConfig) AsIngestionFlow() (*flows.IngestionFlow, error) {
	flow := &flows.IngestionFlow{
		Filetypes: i.Filetypes,
	}
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
		loaderFunc, err := documentloader.GetDocumentLoaderFunc(name, cfg)
		if err != nil {
			return nil, err
		}
		flow.Load = loaderFunc
	}

	if i.TextSplitter.Name != "" {
		name := strings.ToLower(strings.Trim(i.TextSplitter.Name, " "))
		cfg, err := textsplitter.GetTextSplitterConfig(name)
		if err != nil {
			return nil, err
		}
		if len(i.TextSplitter.Options) > 0 {
			jsondata, err := json.Marshal(i.TextSplitter.Options)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(jsondata, &cfg)
			if err != nil {
				return nil, err
			}
		}
		splitterFunc, err := textsplitter.GetTextSplitterFunc(name, cfg)
		if err != nil {
			return nil, err
		}
		flow.Split = splitterFunc
	}

	// TODO: Transformers

	return flow, nil
}

func (f *FlowConfig) ForDataset(name string) (*FlowConfigEntry, error) {
	flowref, ok := f.Datasets[name]
	if ok {
		return f.GetFlow(flowref)
	}
	return f.GetDefaultFlowConfigEntry()
}
