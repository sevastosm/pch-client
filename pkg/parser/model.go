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

type FetchResponse struct {
	Server  IXPServerOption
	Summary domain.BGPSummary
	Nonce   Nonce
}

type Nonce string
