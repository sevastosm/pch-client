package parser

import (
	"github.com/sermojohn/postgres-client/pkg/config"
	"log"
)

// filterServers applies all provided filter config to select list of IXP server options
func filterServers(servers IXPServerOptions, conf config.ParserConfig) []IXPServerOption {
	if ixp := conf.IXP; len(ixp) > 0 {
		servers = servers.filterBy(func(opt *IXPServerOption) bool {
			return config.SliceContains(ixp, opt.IXP)
		})
	}

	if city := conf.City; len(city) > 0 {
		servers = servers.filterBy(func(opt *IXPServerOption) bool {
			return config.SliceContains(city, opt.City)
		})
	}

	if country := conf.Country; len(country) > 0 {
		servers = servers.filterBy(func(opt *IXPServerOption) bool {
			return config.SliceContains(country, opt.Country)
		})
	}

	if limit := conf.ServerLimit; limit > 0 {
		if limit < len(servers) {
			servers = servers[:conf.ServerLimit]
		}
	}

	log.Printf("[parser] filtered IXP servers to size %d\n", len(servers))
	return servers
}
