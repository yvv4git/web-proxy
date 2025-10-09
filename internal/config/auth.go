package config

type Auth struct {
	PredifinedAuth PredifinedAuth `toml:"predifined_auth"`
	// TODO: add more auth types
}

type PredifinedAuth struct {
	Accounts []Account `toml:"accounts"`
}

type Account struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
}
