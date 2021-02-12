package config

type ClientConfig struct {
	IXP         string
	City        string
	Country     string
	IPVersion   string
	ServerLimit int
	ParserRateLimitDelayMillis int64
}
