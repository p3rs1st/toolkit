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

func loadConfig(o *RootOptions, cmd *cobra.Command) func() {
	return func() {
		o.viper = viper.New()
		if info, err := os.Stat(o.ConfigFilepath); err != nil {
			if err := os.MkdirAll(filepath.Dir(o.ConfigFilepath), 0755); err != nil {
				cmd.PrintErrf("Error creating config directory %q\n  %v\n", filepath.Dir(o.ConfigFilepath), err)
				os.Exit(1)
			}
			if _, err := os.OpenFile(o.ConfigFilepath, os.O_RDONLY|os.O_CREATE, 0755); err != nil {
				cmd.PrintErrf("Error opening config file %q\n  %v\n", o.ConfigFilepath, err)
				os.Exit(1)
			}
		} else if info.IsDir() {
			cmd.PrintErrf("Config file %q is a directory\n", o.ConfigFilepath)
			os.Exit(1)
		}
		o.viper.SetConfigFile(o.ConfigFilepath)
		o.viper.SetDefault("current-context", defaultCurrentContext)
		if err := o.viper.ReadInConfig(); err != nil {
			cmd.PrintErrf("Error reading config file %q\n  %v\n", o.ConfigFilepath, err)
			os.Exit(1)
		}
		if err := o.viper.Unmarshal(&o.config); err != nil {
			cmd.PrintErrf("Error parsing config file %q\n  %v\n", o.ConfigFilepath, err)
			os.Exit(1)
		}
		if o.config.CurrentContext == "" {
			o.config.CurrentContext = defaultCurrentContext
		}
		if o.CurrentContext != "" {
			o.config.CurrentContext = o.CurrentContext
		}
		ctx := o.config.GetCurrentContext()
		if o.config.Contexts.GetByName(o.config.CurrentContext) == nil {
			o.config.Contexts = append(o.config.Contexts, ctx)
		}
	}
}

func (o *RootOptions) GetConfig(cmd *cobra.Command) ConfigContext {
	o.once.Do(loadConfig(o, cmd))
	return o.config.GetCurrentContext()
}

func (o *RootOptions) GetRawConfig(cmd *cobra.Command) Config {
	o.once.Do(loadConfig(o, cmd))
	return o.config
}

func (o *RootOptions) SaveCurrentContext(cmd *cobra.Command, name string) {
	config := o.GetRawConfig(cmd)
	if config.Contexts.GetByName(name) == nil {
		o.MergeSaveConfigContext(cmd, config.GetCurrentContext())
	}
	if err := o.viper.MergeConfigMap(map[string]any{"current-context": name}); err != nil {
		cmd.PrintErrf("Error merging config\n  %v\n", err)
		os.Exit(1)
	}
	if err := o.viper.WriteConfig(); err != nil {
		cmd.PrintErrf("Error writing config file %q\n  %v\n", o.ConfigFilepath, err)
		os.Exit(1)
	}
}

func (o *RootOptions) MergeSaveConfigContext(cmd *cobra.Command, configContext ConfigContext) {
	config := o.GetRawConfig(cmd)
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
	if err := o.viper.MergeConfigMap(confMap); err != nil {
		cmd.PrintErrf("Error merging config\n  %v\n", err)
		os.Exit(1)
	}
	if err := o.viper.WriteConfig(); err != nil {
		cmd.PrintErrf("Error writing config file %q\n  %v\n", o.ConfigFilepath, err)
		os.Exit(1)
	}
}

func (o *RootOptions) AllSettings(cmd *cobra.Command) map[string]any {
	o.once.Do(loadConfig(o, cmd))
	settings := o.viper.AllSettings()
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
