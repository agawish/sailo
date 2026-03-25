package commands

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/agawish/sailo/pkg/config"
	"github.com/agawish/sailo/pkg/detect"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize sailo in the current project",
	Long: `Scans the current project for existing Docker configuration
(Dockerfile, docker-compose.yml, devcontainer.json) and creates
a .sailo.yaml config file with sensible defaults.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")
		logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelWarn}))
		return runInit(cmd.OutOrStdout(), ".", force, logger)
	},
}

func init() {
	initCmd.Flags().Bool("force", false, "Overwrite existing .sailo.yaml")
	rootCmd.AddCommand(initCmd)
}

func runInit(out io.Writer, dir string, force bool, logger *slog.Logger) error {
	// Check if config already exists
	existing, err := config.LoadProjectConfig(dir)
	if err != nil {
		return err
	}
	if existing != nil && !force {
		fmt.Fprintln(out, ".sailo.yaml already exists. Use --force to overwrite.")
		return nil
	}

	// Run project detection
	d := detect.NewDetector(logger)
	result, err := d.Detect(dir)
	if err != nil {
		return fmt.Errorf("detect project: %w", err)
	}

	// Build config from detection results
	cfg := config.ProjectConfig{
		Version: 1,
	}

	// Only set base image if no Dockerfile was detected
	if result.Dockerfile == "" && result.BaseImage != "" {
		cfg.Image = result.BaseImage
	}

	// Convert detected ports to port config
	if len(result.Ports) > 0 {
		cfg.Ports = make(map[int]string, len(result.Ports))
		for _, p := range result.Ports {
			cfg.Ports[p] = "auto"
		}
	}

	// Write config
	if err := config.SaveProjectConfig(dir, &cfg); err != nil {
		return err
	}

	// Print summary
	fmt.Fprintln(out, "Initialized .sailo.yaml")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "Detected:")

	if result.Language != "" {
		fmt.Fprintf(out, "  Language:   %s\n", result.Language)
	}
	if result.BaseImage != "" {
		fmt.Fprintf(out, "  Base image: %s\n", result.BaseImage)
	}
	if result.Dockerfile != "" {
		fmt.Fprintln(out, "  Dockerfile: found (will be reused)")
	}
	if result.DockerCompose != "" {
		fmt.Fprintln(out, "  Compose:    found")
	}
	if result.DevContainer != "" {
		fmt.Fprintln(out, "  Devcontainer: found")
	}
	fmt.Fprintf(out, "  Ports:      %s\n", detect.PortSummary(result.Ports))
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "Edit .sailo.yaml to customize your workspace configuration.")

	return nil
}
