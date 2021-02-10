package config

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

func (dbc *DBConfig) HasMissingConfig() bool {
	return dbc.Host == "" || dbc.User == "" || dbc.Password == "" || dbc.User == "" || dbc.Port == 0
}
