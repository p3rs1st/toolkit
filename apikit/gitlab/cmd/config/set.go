package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"slices"
	"strings"
	"toolkit/apikit/gitlab/internal/api"
	"toolkit/apikit/gitlab/internal/types"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func NewSetCommand(op *types.RootOptions) *cobra.Command {
	showKeys := false
	keys := []string{"base_url", "token"}

	cmd := &cobra.Command{
		Use:          "set key value",
		Short:        "Set specified config to configuration file",
		GroupID:      "config sub",
		SilenceUsage: true,
		Args: func(cmd *cobra.Command, args []string) error {
			if showKeys {
				return nil
			}
			if len(args) < 1 {
				return errors.New("requires at least one argument")
			}
			if args[0] == "token" {
				if len(args) > 1 {
					return errors.New("Too many arguments for token")
				}
				return nil
			}
			if err := cobra.ExactArgs(2)(cmd, args); err != nil {
				return err
			}
			if !slices.Contains(keys, args[0]) {
				return errors.New("Invalid config key: " + args[0])
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if showKeys {
				cmd.Printf("config keys:\n%s\n", strings.Join(keys, " "))
				return nil
			}

			conf := op.GetConfig(cmd)
			args[0] = strings.ToLower(args[0])
			if args[0] == "base_url" {
				uri, err := url.ParseRequestURI(args[1])
				if err != nil || (uri.Scheme != "http" && uri.Scheme != "https") || uri.Host == "" {
					return fmt.Errorf("%q is not a valid URL", args[1])
				}
				conf.BaseURL = args[1]
			} else if args[0] == "token" {
				cmd.Print("Enter your token: ")
				tokenBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
				cmd.Println()
				if err != nil {
					return fmt.Errorf("term.ReadPassword failed with, %s\n", err.Error())
				}

				args = append(args, string(tokenBytes))
				conf.Token = args[1]
				if ok, err := api.CheckAccessToken(conf); !ok {
					if err == nil {
						err = api.ErrNoAuthorization
					}
					return err
				}
			}
			op.MergeSaveConfigContext(cmd, conf)

			return nil
		},
	}

	cmd.Flags().BoolVar(&showKeys, "show", false, "Show all keys in configuration")

	return cmd
}
