package configuration

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strconv"
	"strings"
)

var (

	// ErrNoInstance is returned when an instance definition is expected but missing
	ErrNoInstance = errors.New("missing instance definition")

	// ErrInvalidInstanceDefinition is returned when an invalid instance definition is discovered, e.g. an empty name or address
	ErrInvalidInstanceDefinition = errors.New("invalid instance definition")

	// ErrDuplicateInstance is returned when there are multiple definitions for the same instance
	ErrDuplicateInstance = errors.New("duplicate instance")
)

// Instance describes a single  instance connection information
type Instance struct {
	Name    string `yaml:"name"`
	Address string `yaml:"address"`
}

// InstanceConfig describes the set of peers and clients in the system
type InstanceConfig struct {
	Peers   []Instance `yaml:"peers"`
	Clients []Instance `yaml:"clients"`
}

// NewInstanceConfig loads a  instance configuration from given file
func NewInstanceConfig(fname string, name int64) (*InstanceConfig, error) {
	var cfg InstanceConfig

	data, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	err = yaml.UnmarshalStrict(data, &cfg)
	if err != nil {
		return nil, err
	}
	cfg = configureSelfIP(cfg, name)
	// sanity checks
	if err := checkInstanceList(cfg.Peers...); err != nil {
		return nil, err
	}
	if err := checkInstanceList(cfg.Clients...); err != nil {
		return nil, err
	}
	return &cfg, nil
}

/*
	Replace the IP part of my self to 0.0.0.0
*/

func configureSelfIP(cfg InstanceConfig, name int64) InstanceConfig {
	for i := 0; i < len(cfg.Peers); i++ {
		if cfg.Peers[i].Name == strconv.Itoa(int(name)) {
			cfg.Peers[i].Address = "0.0.0.0:" + getPort(cfg.Peers[i].Address)
			return cfg
		}
	}
	for i := 0; i < len(cfg.Clients); i++ {
		if cfg.Clients[i].Name == strconv.Itoa(int(name)) {
			cfg.Clients[i].Address = "0.0.0.0:" + getPort(cfg.Clients[i].Address)
			return cfg
		}
	}
	return cfg
}

/*
	Returns the port part of the ip:port
*/

func getPort(address string) string {
	return strings.Split(address, ":")[1]
}

// checkInstanceList performs a sanity check for a list of instances
// this method checks whether the identifier is unique and the address is unique for each replica and client
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
