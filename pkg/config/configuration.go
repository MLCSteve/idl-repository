package config

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Repository   string       `yaml:"repository"`
	Name         string       `yaml:"name"`
	IdlDirectory string       `yaml:"idl_directory"`
	Dependencies []Dependency `yaml:"dependencies,omitempty"`
	Provides     []Provide    `yaml:"provides,omitempty"`
}

type Dependency struct {
	Name       string `yaml:"name"`
	Version    string `yaml:"version"`
	Type       string `yaml:"type"`
	Repository string `yaml:"repository"`
}

type Provide struct {
	Root      string `yaml:"root"`
	Type      string `yaml:"type"`
	IdlIgnore string `yaml:"idlignore"`
}

func (c *Configuration) Marshal(writer io.Writer) error {
	en := yaml.NewEncoder(writer)
	defer en.Close()

	err := en.Encode(c)
	if err != nil {
		return errors.Wrap(err, "failed to encode configuration")
	}

	return nil
}

func (c *Configuration) UnMarshal(reader io.Reader) error {
	d := yaml.NewDecoder(reader)

	err := d.Decode(c)
	if err != nil {
		return errors.Wrap(err, "failed to decode configuration")
	}

	return nil
}

func (c *Configuration) ResolveRepository(dependency Dependency) string {
	if dependency.Repository != "" {
		return dependency.Repository
	}

	return c.Repository
}

func (c *Configuration) Validate() error {
	requires := map[string]map[string]bool{}
	for _, dep := range c.Dependencies {
		_, ok := requires[dep.Name]
		if !ok {
			requires[dep.Name] = map[string]bool{}
		}
		_, ok = requires[dep.Name][dep.Type]
		if ok {
			return errors.New(fmt.Sprintf("the dependency '%s' with type '%s' has more than one entry", dep.Name, dep.Type))
		}
		requires[dep.Name][dep.Type] = true
	}
	return nil
}
