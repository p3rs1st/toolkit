package types

type Config struct {
	CurrentContext string         `mapstructure:"current-context" json:"current-context"`
	Contexts       ConfigContexts `mapstructure:"contexts" json:"contexts"`
}

func (c Config) GetCurrentContext() ConfigContext {
	ctx := c.Contexts.GetByName(c.CurrentContext)
	if ctx == nil {
		return ConfigContext{
			Name:    c.CurrentContext,
			BaseURL: "https://gitlab.com",
			Token:   "",
		}
	}
	return *ctx
}

type ConfigContext struct {
	Name    string `mapstructure:"name" json:"name"`
	BaseURL string `mapstructure:"base_url" json:"base_url"`
	Token   string `mapstructure:"token" json:"token"`
}

type ConfigContexts []ConfigContext

func (c ConfigContexts) GetByName(name string) *ConfigContext {
	idx := c.GetIdxByName(name)
	if idx == -1 {
		return nil
	}
	return &c[idx]
}

func (c ConfigContexts) GetIdxByName(name string) int {
	if name == "" {
		return -1
	}
	for i, context := range c {
		if context.Name == name {
			return i
		}
	}
	return -1
}
