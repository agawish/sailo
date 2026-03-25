package commands

import (
	"fmt"
	"io"
	"strings"

	"github.com/agawish/sailo/pkg/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage sailo configuration",
	Long: `View and modify sailo configuration.

Configuration files:
  Project: .sailo.yaml (checked into repo)
  User:    ~/.sailo/config.yaml

Example:
  sailo config show
  sailo config set defaults.from develop
  sailo config set defaults.port_range 4001-4999`,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runConfigShow(cmd.OutOrStdout())
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runConfigSet(cmd.OutOrStdout(), args[0], args[1])
	},
}

func init() {
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
	rootCmd.AddCommand(configCmd)
}

func runConfigShow(out io.Writer) error {
	userCfg, err := config.LoadUserConfig()
	if err != nil {
		return fmt.Errorf("load user config: %w", err)
	}

	fmt.Fprintln(out, "User config (~/.sailo/config.yaml):")
	fmt.Fprintf(out, "  defaults.from:          %s\n", userCfg.Defaults.From)
	fmt.Fprintf(out, "  defaults.cleanup_after: %s\n", userCfg.Defaults.CleanupAfter)
	fmt.Fprintf(out, "  defaults.port_range:    %s\n", userCfg.Defaults.PortRange)
	fmt.Fprintf(out, "  env_passthrough:        %s\n", strings.Join(userCfg.EnvPassthrough, ", "))
	fmt.Fprintf(out, "  git.credentials:        %s\n", userCfg.Git.Credentials)
	fmt.Fprintf(out, "  git.auto_push:          %v\n", userCfg.Git.AutoPush)

	projectCfg, err := config.LoadProjectConfig(".")
	if err != nil {
		return fmt.Errorf("load project config: %w", err)
	}

	fmt.Fprintln(out, "")
	if projectCfg == nil {
		fmt.Fprintln(out, "No .sailo.yaml found in current directory.")
	} else {
		fmt.Fprintln(out, "Project config (.sailo.yaml):")
		fmt.Fprintf(out, "  version: %d\n", projectCfg.Version)
		if projectCfg.Image != "" {
			fmt.Fprintf(out, "  image:   %s\n", projectCfg.Image)
		}
		if len(projectCfg.Ports) > 0 {
			for port, mapping := range projectCfg.Ports {
				fmt.Fprintf(out, "  port:    %d → %s\n", port, mapping)
			}
		}
		if len(projectCfg.Services) > 0 {
			fmt.Fprintf(out, "  services: %s\n", strings.Join(projectCfg.Services, ", "))
		}
		if projectCfg.Test != "" {
			fmt.Fprintf(out, "  test:    %s\n", projectCfg.Test)
		}
	}

	return nil
}

// validConfigKeys lists the supported dot-notation keys for config set.
var validConfigKeys = []string{
	"defaults.from",
	"defaults.cleanup_after",
	"defaults.port_range",
	"git.credentials",
	"git.auto_push",
	"env_passthrough",
}

func runConfigSet(out io.Writer, key, value string) error {
	cfg, err := config.LoadUserConfig()
	if err != nil {
		return fmt.Errorf("load user config: %w", err)
	}

	switch key {
	case "defaults.from":
		cfg.Defaults.From = value
	case "defaults.cleanup_after":
		cfg.Defaults.CleanupAfter = value
	case "defaults.port_range":
		cfg.Defaults.PortRange = value
	case "git.credentials":
		cfg.Git.Credentials = value
	case "git.auto_push":
		switch strings.ToLower(value) {
		case "true", "1", "yes":
			cfg.Git.AutoPush = true
		case "false", "0", "no":
			cfg.Git.AutoPush = false
		default:
			return fmt.Errorf("invalid boolean value for git.auto_push: %q (use true/false)", value)
		}
	case "env_passthrough":
		parts := strings.Split(value, ",")
		trimmed := make([]string, 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				trimmed = append(trimmed, p)
			}
		}
		cfg.EnvPassthrough = trimmed
	default:
		return fmt.Errorf("unknown config key: %q\nValid keys: %s", key, strings.Join(validConfigKeys, ", "))
	}

	if err := config.SaveUserConfig(cfg); err != nil {
		return fmt.Errorf("save user config: %w", err)
	}

	fmt.Fprintf(out, "Set %s = %s\n", key, value)
	return nil
}
