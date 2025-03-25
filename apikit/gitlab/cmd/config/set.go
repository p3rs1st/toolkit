package config

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"toolkit/apikit/gitlab/cmd/util/uurl"
	"toolkit/apikit/gitlab/internal/api"
	"toolkit/apikit/gitlab/internal/types"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	errInvalidKey = errors.New("invalid config key")
	errReadPasswd = errors.New("failed with term.ReadPassword")

	setCommandKeys = []string{"base_url", "token"}
)

func NewSetCommand(option *types.RootOptions) *cobra.Command {
	showKeys := false

	cmd := &cobra.Command{
		Use:          "set key value",
		Short:        "Set specified config to configuration file",
		GroupID:      "config sub",
		SilenceUsage: true,
		Args: func(cmd *cobra.Command, args []string) error {
			if showKeys {
				return nil
			}
			if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
				return err
			}
			if args[0] == "token" {
				return cobra.MaximumNArgs(1)(cmd, args)
			}
			if err := cobra.ExactArgs(2)(cmd, args); err != nil {
				return err
			}
			if !slices.Contains(setCommandKeys, args[0]) {
				return fmt.Errorf("%w: %s", errInvalidKey, args[0])
			}

			return nil
		},
		RunE: setCommandRunE(option, &showKeys),
	}

	cmd.Flags().BoolVar(&showKeys, "show", false, "Show all keys in configuration")

	return cmd
}

func setCommandRunE(
	option *types.RootOptions, showKeys *bool,
) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if *showKeys {
			cmd.Printf("config keys:\n%s\n", strings.Join(setCommandKeys, " "))

			return nil
		}

		conf := option.GetConfig(cmd)
		args[0] = strings.ToLower(args[0])

		if args[0] == "base_url" {
			if err := uurl.CheckURLValid(args[1]); err != nil {
				return err
			}

			conf.BaseURL = args[1]
		}

		if args[0] == "token" {
			cmd.Print("Enter your token: ")

			tokenBytes, err := term.ReadPassword(int(os.Stdin.Fd()))

			cmd.Println()

			if err != nil {
				return fmt.Errorf("%w: %w", errReadPasswd, err)
			}

			args = append(args, string(tokenBytes))
			conf.Token = args[1]

			ok, err := api.CheckAccessToken(conf)
			if err != nil {
				return err
			}

			if !ok {
				return api.ErrNoAuthorization
			}
		}

		option.MergeSaveConfigContext(cmd, conf)

		return nil
	}
}
