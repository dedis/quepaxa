package configuration

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strconv"
)

/*
	config.go define the structs and methods to pass the configuration file, that contains the IP:port of each client and replica
*/

// ReplicaInstance describes a single  QuePaxa connection information
type ReplicaInstance struct {
	Name         string `yaml:"name"` // unique id of the QuePaxa node
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

// InstanceConfig describes the set of replicas and clients in the system
type InstanceConfig struct {
	Peers   []ReplicaInstance `yaml:"peers"`
	Clients []ClientInstance  `yaml:"clients"`
}

// NewInstanceConfig loads an instance configuration from given file
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
	// set the self ip to 0.0.0.0
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
