package configuration

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strconv"
)

// ReplicaInstance describes a single  instance connection information
type ReplicaInstance struct {
	Name         string `yaml:"name"`
	IP           string `yaml:"ip"`
	PROXYPORT    string `yaml:"proxyport"`
	RECORDERPORT string `yaml:"recorderport"`
}

// ClientInstance describes a single CLIENT connection information
type ClientInstance struct {
	Name       string `yaml:"name"`
	IP         string `yaml:"ip"`
	CLIENTPORT string `yaml:"clientport"`
}

// InstanceConfig describes the set of peers and clients in the system
type InstanceConfig struct {
	Peers   []ReplicaInstance `yaml:"peers"`
	Clients []ClientInstance  `yaml:"clients"`
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
	return &cfg, nil
}

/*
	Replace the IP of my self to 0.0.0.0
*/

func configureSelfIP(cfg InstanceConfig, name int64) InstanceConfig {
	for i := 0; i < len(cfg.Peers); i++ {
		if cfg.Peers[i].Name == strconv.Itoa(int(name)) {
			cfg.Peers[i].IP = "0.0.0.0"
			return cfg
		}
	}
	for i := 0; i < len(cfg.Clients); i++ {
		if cfg.Clients[i].Name == strconv.Itoa(int(name)) {
			cfg.Clients[i].IP = "0.0.0.0"
			return cfg
		}
	}
	return cfg
}
