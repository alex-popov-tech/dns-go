package config

import "net/netip"

type Config struct {
	Resolver netip.AddrPort
}
