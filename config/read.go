package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

// Read reads configuration from file
func Read(name string) (*Configuration, error) {
	if len(name) == 0 {
		return nil, errors.New("no configuration file name given")
	}

	// Reading configuration file
	bts, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}

	// Parsing yaml
	var target Configuration
	if err := yaml.Unmarshal(bts, &target); err != nil {
		return nil, err
	}

	// Validating
	if err := target.Validate(); err != nil {
		return nil, err
	}

	return &target, nil
}
