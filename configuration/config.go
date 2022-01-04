package configuration

import (
	"errors"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
)

var (

	// ErrNoInstance is returned when an instance definition is expected but missing
	ErrNoInstance = errors.New("missing instance definition")

	// ErrInvalidInstanceDefinition is returned when an invalid instance definition is discovered, e.g. an empty name or
	// address
	ErrInvalidInstanceDefinition = errors.New("invalid instance definition")

	// ErrDuplicateInstance is returned when there are multiple definitions for the same instance
	ErrDuplicateInstance = errors.New("duplicate instance")
)

// Instance describes a single  instance connection information
type Instance struct {
	Name    string `yaml:"name"`
	Address string `yaml:"address"`
}

// InstanceConfig describes a  instance configuration
type InstanceConfig struct {
	Name   string     `yaml:"name"`
	Listen string     `yaml:"listen"`
	Peers  []Instance `yaml:"peers"`
}

// QuorumConfig describes a  quorum configuration file
type QuorumConfig struct {
	Instances []Instance `yaml:"instances"`
}

// NewInstanceConfig loads a  instance configuration from given file
func NewInstanceConfig(fname string) (*InstanceConfig, error) {
	var cfg InstanceConfig

	data, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	err = yaml.UnmarshalStrict(data, &cfg)
	if err != nil {
		return nil, err
	}

	// sanity checks
	instances := append(cfg.Peers, Instance{Name: cfg.Name, Address: cfg.Listen})
	if err := checkInstanceList(instances...); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// NewQuorumConfig loads a  quorum configuration from given file
func NewQuorumConfig(fname string) (*QuorumConfig, error) {
	cfg := &QuorumConfig{}
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	err = yaml.UnmarshalStrict(data, cfg)
	if err != nil {
		return nil, err
	}

	// sanity checks
	if err := checkInstanceList(cfg.Instances...); err != nil {
		return nil, err
	}

	return cfg, nil
}

// checkInstanceList performs a sanity check for a list of instances
func checkInstanceList(instances ...Instance) error {
	if len(instances) == 0 {
		return ErrNoInstance
	}

	seenNames := make(map[string]bool)
	seenAddresses := make(map[string]bool)
	for _, in := range instances {
		// Instance name must be unique in the configuration file.
		if len(in.Name) == 0 {
			return ErrInvalidInstanceDefinition
		}
		if seenNames[in.Name] {
			return ErrDuplicateInstance
		}
		seenNames[in.Name] = true

		// Instance address must be unique in the configuration file.
		if len(in.Address) == 0 {
			return ErrInvalidInstanceDefinition
		}
		if seenAddresses[in.Address] {
			return ErrDuplicateInstance
		}
		seenAddresses[in.Address] = true
	}
	return nil
}
