package types

type Config struct {
	CurrentContext string         `json:"current-context" mapstructure:"current-context"`
	Contexts       ConfigContexts `json:"contexts"        mapstructure:"contexts"`
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
	Name    string `json:"name"     mapstructure:"name"`
	BaseURL string `json:"base_url" mapstructure:"base_url"`
	Token   string `json:"token"    mapstructure:"token"`
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
