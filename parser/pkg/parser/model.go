package parser

import "github.com/sermojohn/postgres-client/pkg/domain"

type InitResponse struct {
	Servers []IXPServerOption
	Nonce   Nonce
}

type IXPServerOption struct {
	domain.IXPServer
	ItemID int
}

type IXPServerOptions []IXPServerOption

func (opts IXPServerOptions) filterBy(filter func(option *IXPServerOption) bool) IXPServerOptions {
	var out IXPServerOptions

	for _, opt := range opts {
		if filter(&opt) {
			out = append(out, opt)
		}
	}

	return out
}

type FetchResponse struct {
	Server  IXPServerOption
	Summary domain.BGPSummary
	Nonce   Nonce
}

type Nonce string
