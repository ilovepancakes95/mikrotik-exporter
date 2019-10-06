package config

import (
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config represents the configuration for the exporter
type Config struct {
	Devices  []Device `yaml:"devices"`
	Features struct {
		BGP            bool `yaml:"bgp,omitempty"`
		DHCP           bool `yaml:"dhcp,omitempty"`
		DHCPLeases     bool `yaml:"dhcp-leases,omitempty"`
		DHCPv6         bool `yaml:"dhcpv6,omitempty"`
		Routes         bool `yaml:"routes,omitempty"`
		RoutesV6       bool `yaml:"routesv6,omitempty"`
		Pool           bool `yaml:"pool,omitempty"`
		PoolV6         bool `yaml:"poolv6,omitempty"`
		Optics         bool `yaml:"optics,omitempty"`
		WlanStations   bool `yaml:"wlan-stations,omitempty"`
		WlanInterfaces bool `yaml:"wlan-interfaces,omitempty"`
		Monitor        bool `yaml:"monitor,omitempty"`
		IPSecPeers     bool `yaml:"ipsec-peers,omitempty"`
		OSPFNeighbor   bool `yaml:"ospf-neighbor,omitempty"`
	} `yaml:"features,omitempty"`
}

// Device represents a target device
type Device struct {
	Name     string `yaml:"name"`
	Address  string `yaml:"address"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// Load reads YAML from reader and unmashals in Config
func Load(r io.Reader) (*Config, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = yaml.Unmarshal(b, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
