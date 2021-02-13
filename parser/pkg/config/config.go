package config

type ClientConfig struct {
	IXP         string
	City        string
	Country     string
	IPVersion   string // IPv4 or IPv6
	ServerLimit int
	ParserRateLimitDelayMillis int64
}

const (
	IPv4 = "IPv4"
	IPv6 = "IPv6"
)
