package types

type Config struct {
	BaseURL string `mapstructure:"base_url"`
	Token   string `mapstructure:"token"`
}
