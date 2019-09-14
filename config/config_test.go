package config

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestShouldParse(t *testing.T) {
	b := loadTestFile(t)
	c, err := Load(bytes.NewReader(b))
	if err != nil {
		t.Fatalf("could not parse: %v", err)
	}

	if len(c.Devices) != 2 {
		t.Fatalf("expected 2 devices, got %v", len(c.Devices))
	}

	assertDevice("test1", "192.168.1.1", "foo", "bar", c.Devices[0], t)
	assertDevice("test2", "192.168.2.1", "test", "123", c.Devices[1], t)
	assertFeature("BGP", c.Features.BGP, t)
	assertFeature("DHCP", c.Features.DHCP, t)
	assertNoFeature("DHCPv6", c.Features.DHCPv6, t)
	assertFeature("Pool", c.Features.Pool, t)
	assertNoFeature("PoolV6", c.Features.PoolV6, t)
	assertFeature("Routes", c.Features.Routes, t)
	assertNoFeature("RoutesV6", c.Features.RoutesV6, t)
	assertFeature("Optics", c.Features.Optics, t)
	assertFeature("WlanStations", c.Features.WlanStations, t)
	assertFeature("WlanInterfaces", c.Features.WlanInterfaces, t)
}

func loadTestFile(t *testing.T) []byte {
	b, err := ioutil.ReadFile("config.test.yml")
	if err != nil {
		t.Fatalf("could not load config: %v", err)
	}

	return b
}

func assertDevice(name, address, user, password string, c Device, t *testing.T) {
	if c.Name != name {
		t.Fatalf("expected name %s, got %s", name, c.Name)
	}

	if c.Address != address {
		t.Fatalf("expected address %s, got %s", address, c.Address)
	}

	if c.User != user {
		t.Fatalf("expected user %s, got %s", user, c.User)
	}

	if c.Password != password {
		t.Fatalf("expected password %s, got %s", password, c.Password)
	}
}

func assertFeature(name string, v bool, t *testing.T) {
	if !v {
		t.Fatalf("expected feature %s to be enabled", name)
	}
}

func assertNoFeature(name string, v bool, t *testing.T) {
	if v {
		t.Fatalf("expected feature %s to be disabled", name)
	}
}
