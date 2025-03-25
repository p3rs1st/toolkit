package types

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"toolkit/apikit/gitlab/pkg"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	HomeDir, _        = os.UserHomeDir()
	DefaultToolkitDir = filepath.Join(HomeDir, ".apikit")
	DefaultConfigPath = filepath.Join(DefaultToolkitDir, "gitlab.yaml")
)

type RootOptions struct {
	ConfigFilepath string
	CurrentContext string

	viper  *viper.Viper
	once   sync.Once
	config Config
}

const (
	defaultCurrentContext string = "default"
)

func loadConfig(option *RootOptions, cmd *cobra.Command) func() {
	return func() {
		option.viper = viper.New()
		if info, err := os.Stat(option.ConfigFilepath); err != nil {
			if err := os.MkdirAll(filepath.Dir(option.ConfigFilepath), 0o755); err != nil {
				cmd.PrintErrf("Error creating config directory %q\n  %v\n", filepath.Dir(option.ConfigFilepath), err)
				os.Exit(1)
			}

			if _, err := os.OpenFile(option.ConfigFilepath, os.O_RDONLY|os.O_CREATE, 0o755); err != nil {
				cmd.PrintErrf("Error opening config file %q\n  %v\n", option.ConfigFilepath, err)
				os.Exit(1)
			}
		} else if info.IsDir() {
			cmd.PrintErrf("Config file %q is a directory\n", option.ConfigFilepath)
			os.Exit(1)
		}

		option.viper.SetConfigFile(option.ConfigFilepath)
		option.viper.SetDefault("current-context", defaultCurrentContext)

		if err := option.viper.ReadInConfig(); err != nil {
			cmd.PrintErrf("Error reading config file %q\n  %v\n", option.ConfigFilepath, err)
			os.Exit(1)
		}

		if err := option.viper.Unmarshal(&option.config); err != nil {
			cmd.PrintErrf("Error parsing config file %q\n  %v\n", option.ConfigFilepath, err)
			os.Exit(1)
		}

		if option.config.CurrentContext == "" {
			option.config.CurrentContext = defaultCurrentContext
		}

		if option.CurrentContext != "" {
			option.config.CurrentContext = option.CurrentContext
		}

		ctx := option.config.GetCurrentContext()
		if option.config.Contexts.GetByName(option.config.CurrentContext) == nil {
			option.config.Contexts = append(option.config.Contexts, ctx)
		}
	}
}

func (option *RootOptions) GetConfig(cmd *cobra.Command) ConfigContext {
	option.once.Do(loadConfig(option, cmd))

	return option.config.GetCurrentContext()
}

func (option *RootOptions) GetRawConfig(cmd *cobra.Command) Config {
	option.once.Do(loadConfig(option, cmd))

	return option.config
}

func (option *RootOptions) SaveCurrentContext(cmd *cobra.Command, name string) {
	config := option.GetRawConfig(cmd)
	if config.Contexts.GetByName(name) == nil {
		option.MergeSaveConfigContext(cmd, config.GetCurrentContext())
	}

	if err := option.viper.MergeConfigMap(map[string]any{"current-context": name}); err != nil {
		cmd.PrintErrf("Error merging config\n  %v\n", err)
		os.Exit(1)
	}

	if err := option.viper.WriteConfig(); err != nil {
		cmd.PrintErrf("Error writing config file %q\n  %v\n", option.ConfigFilepath, err)
		os.Exit(1)
	}
}

func (option *RootOptions) MergeSaveConfigContext(cmd *cobra.Command, configContext ConfigContext) {
	config := option.GetRawConfig(cmd)
	ctxs := config.Contexts

	idx := ctxs.GetIdxByName(config.CurrentContext)
	if idx == -1 {
		config.Contexts = append(config.Contexts, configContext)
	} else {
		ctxs[idx] = configContext
	}

	confMap := map[string]any{
		"contexts": pkg.MapFunc(ctxs, func(ctx ConfigContext) map[string]any {
			ctxBytes, err := json.Marshal(ctx)
			if err != nil {
				cmd.PrintErrf("Error marshaling context %q\n  %v\n", ctx.Name, err)
				os.Exit(1)
			}
			ctxMap := map[string]any{}
			if err := json.Unmarshal(ctxBytes, &ctxMap); err != nil {
				cmd.PrintErrf("Error unmarshaling context %q\n  %v\n", ctx.Name, err)
				os.Exit(1)
			}

			return ctxMap
		}),
	}

	if err := option.viper.MergeConfigMap(confMap); err != nil {
		cmd.PrintErrf("Error merging config\n  %v\n", err)
		os.Exit(1)
	}

	if err := option.viper.WriteConfig(); err != nil {
		cmd.PrintErrf("Error writing config file %q\n  %v\n", option.ConfigFilepath, err)
		os.Exit(1)
	}
}

func (option *RootOptions) AllSettings(cmd *cobra.Command) map[string]any {
	option.once.Do(loadConfig(option, cmd))
	settings := option.viper.AllSettings()
	ctxInterface := settings["contexts"].([]interface{})
	ctxs := make([]map[string]any, len(ctxInterface))

	for i, ctx := range ctxInterface {
		ctxMap, _ := ctx.(map[string]any)
		ctxs[i] = ctxMap
	}

	for i, ctx := range ctxs {
		if token, ok := ctx["token"].(string); ok && token != "" {
			ctxs[i]["token"] = "*****"
		}
	}

	return settings
}
