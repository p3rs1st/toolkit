package types

import (
	"os"
	"path/filepath"
	"sync"

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

	viper  *viper.Viper
	once   sync.Once
	config Config
}

func (o *RootOptions) GetConfig(cmd *cobra.Command) Config {
	o.once.Do(func() {
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
		o.viper.SetDefault("base_url", "https://gitlab.com")
		if err := o.viper.ReadInConfig(); err != nil {
			cmd.PrintErrf("Error reading config file %q\n  %v\n", o.ConfigFilepath, err)
			os.Exit(1)
		}
		if err := o.viper.Unmarshal(&o.config); err != nil {
			cmd.PrintErrf("Error parsing config file %q\n  %v\n", o.ConfigFilepath, err)
			os.Exit(1)
		}
	})
	return o.config
}

func (o *RootOptions) MergeSaveConfig(cmd *cobra.Command, confMap map[string]any) {
	_ = o.GetConfig(cmd)
	if err := o.viper.MergeConfigMap(confMap); err != nil {
		cmd.PrintErrf("Error merging config\n  %v\n", err)
		os.Exit(1)
	}
	if err := o.viper.WriteConfig(); err != nil {
		cmd.PrintErrf("Error writing config file %q\n  %v\n", o.ConfigFilepath, err)
		os.Exit(1)
	}
}

func (o *RootOptions) AllKeys(cmd *cobra.Command) []string {
	_ = o.GetConfig(cmd)
	return o.viper.AllKeys()
}

func (o *RootOptions) AllSettings(cmd *cobra.Command) map[string]any {
	_ = o.GetConfig(cmd)
	return o.viper.AllSettings()
}
