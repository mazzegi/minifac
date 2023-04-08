package minifac

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/mazzegi/minifac/assets"
	"golang.org/x/exp/slices"
)

type Config struct {
	BaseResources []Resource     `json:"base-resources"`
	Resources     []Resource     `json:"resources"`
	Obstacles     []ObstacleType `json:"obstacles"`
	Receipts      []Receipt      `json:"receipts"`
}

func LoadConfigFromFile(file string) (*Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("open %q: %w", file, err)
	}
	defer f.Close()
	return LoadConfig(f)
}

func LoadConfig(r io.Reader) (*Config, error) {
	var cfg Config
	err := json.NewDecoder(r).Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("json.decode: %w", err)
	}
	err = cfg.validate()
	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}
	return &cfg, nil
}

func (c *Config) containsResource(res Resource) bool {
	for _, eres := range c.BaseResources {
		if res == eres {
			return true
		}
	}
	for _, eres := range c.Resources {
		if res == eres {
			return true
		}
	}
	return false
}

func (c *Config) validate() error {
	// check receipts
	for _, rec := range c.Receipts {
		if !c.containsResource(rec.Output) {
			return fmt.Errorf("output resource for receipt %q not found", rec.Output)
		}
		for res := range rec.Input {
			if !c.containsResource(res) {
				return fmt.Errorf("input resource %q for receipt %q not found", res, rec.Output)
			}
		}
	}
	return nil
}

func (c *Config) ValidateAssets(as *assets.Assets) error {
	asRess := as.ResourceNames()
	asObsts := as.ObstacleNames()
	for _, res := range c.BaseResources {
		if !slices.Contains(asRess, string(res)) {
			return fmt.Errorf("no resource asset for %q", res)
		}
	}
	for _, res := range c.Resources {
		if !slices.Contains(asRess, string(res)) {
			return fmt.Errorf("no resource asset for %q", res)
		}
	}
	for _, obs := range c.Obstacles {
		if !slices.Contains(asObsts, string(obs)) {
			return fmt.Errorf("no obstacle asset for %q", obs)
		}
	}
	return nil
}
