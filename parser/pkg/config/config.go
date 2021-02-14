package config

import "strings"

type Config struct {
	DBConfig     `yaml:"db"`
	ParserConfig `yaml:"parser"`
}
type ParserConfig struct {
	IXP                        []string `yaml:"ixp"`
	City                       []string `yaml:"city"`
	Country                    []string `yaml:"country"`
	Protocol                   string   `yaml:"protocol"` // IPv4 or IPv6
	ServerLimit                int      `yaml:"server_limit"`
	ParserRateLimitDelayMillis int64    `yaml:"rate_limit_delay_ms"`
}

func SliceContains(ar []string, cand string) bool {
	for _, el := range ar {
		if strings.ToLower(el) == strings.ToLower(cand) {
			return true
		}
	}
	return false
}

const (
	IPv4 = "IPv4"
	IPv6 = "IPv6"
)

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

func (c *Config) Valid() (bool, string) {
	if c.Protocol != IPv4 && c.Protocol != IPv6 {
		return false, "protocol does not have expected value [IPv4, IPv6]"
	}
	if !c.DBConfig.Valid() {
		return false, "DB configuration is missing required values"
	}
	return true, ""
}

func (dbc *DBConfig) Valid() bool {
	return !(dbc.Host == "" || dbc.User == "" || dbc.Password == "" || dbc.Database == "" || dbc.Port == 0)
}
