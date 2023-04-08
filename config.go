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
	Assemblers    map[Resource]struct {
		Input               map[Resource]int `json:"input"`
		InputStockCapacity  int              `json:"input-stock-capacity"`
		OutputStockCapacity int              `json:"output-stock-capacity"`
		ProductionTime      int              `json:"production-time"`
	} `json:"assemblers"`
	Producers map[Resource]struct {
		Rate          Rate `json:"rate"`
		StockCapacity int  `json:"stock-capacity"`
	} `json:"producers"`
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

func (c *Config) containsObstacle(obs ObstacleType) bool {
	for _, eobs := range c.Obstacles {
		if obs == eobs {
			return true
		}
	}
	return false
}

func (c *Config) validate() error {
	if !c.containsObstacle(ObstacleWall) {
		return fmt.Errorf("config contains no wall-obstacle")
	}

	// check receipts
	for ares, rec := range c.Assemblers {
		if !c.containsResource(ares) {
			return fmt.Errorf("output resource for receipt %q not found", ares)
		}
		for res := range rec.Input {
			if !c.containsResource(res) {
				return fmt.Errorf("input resource %q for receipt %q not found", res, ares)
			}
		}
	}
	for res := range c.Producers {
		if !c.containsResource(res) {
			return fmt.Errorf("producer resource %q not found", res)
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

func (c *Config) ValidatePuzzle(pzl *Puzzle) error {
	for _, e := range pzl.Producers {
		if _, ok := c.Producers[e.Resource]; !ok {
			return fmt.Errorf("no producer for resource %q", e.Resource)
		}
	}
	for _, e := range pzl.Assemblers {
		if _, ok := c.Assemblers[e.Resource]; !ok {
			return fmt.Errorf("no assembler for resource %q", e.Resource)
		}
	}
	for _, e := range pzl.Finalizers {
		if !c.containsResource(e.Resource) {
			return fmt.Errorf("finalizer resource %q not available", e.Resource)
		}
	}
	for _, e := range pzl.Obstacles {
		if !c.containsObstacle(e.Type) {
			return fmt.Errorf("obstacle type %q not available", e.Type)
		}
	}
	return nil
}
