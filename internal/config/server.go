package config

type Server struct {
	Host    string `toml:"host"`
	Port    uint16 `toml:"port"`
	Verbose bool   `toml:"verbose"`
}
