package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type dbInitOptions struct {
	projectRoot string
	service     string
	reset       bool
	timeout     time.Duration
}

func main() {
	var opts dbInitOptions

	rootCmd := &cobra.Command{
		Use:   "gomall",
		Short: "GoMall Lite command line tools",
	}

	dbCmd := &cobra.Command{
		Use:   "db",
		Short: "Database commands",
	}

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Run database/init.sql inside the MySQL container",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDBInit(opts)
		},
	}

	initCmd.Flags().StringVar(&opts.projectRoot, "project-root", "", "project root containing docker-compose.yml and database/init.sql")
	initCmd.Flags().StringVar(&opts.service, "service", "mysql", "docker compose MySQL service name")
	initCmd.Flags().BoolVar(&opts.reset, "reset", false, "remove compose containers and volumes before initialization")
	initCmd.Flags().DurationVar(&opts.timeout, "timeout", 2*time.Minute, "maximum time to wait for MySQL healthcheck")

	dbCmd.AddCommand(initCmd)
	rootCmd.AddCommand(dbCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runDBInit(opts dbInitOptions) error {
	projectRoot, err := resolveProjectRoot(opts.projectRoot)
	if err != nil {
		return err
	}

	sqlPath := filepath.Join(projectRoot, "database", "init.sql")
	if _, err := os.Stat(sqlPath); err != nil {
		return fmt.Errorf("database init sql not found: %w", err)
	}

	if opts.reset {
		fmt.Println("Resetting containers and volumes...")
		if err := run(projectRoot, "docker", "compose", "down", "-v"); err != nil {
			return err
		}
	}

	fmt.Println("Starting MySQL container...")
	if err := run(projectRoot, "docker", "compose", "up", "-d", opts.service); err != nil {
		return err
	}

	containerID, err := output(projectRoot, "docker", "compose", "ps", "-q", opts.service)
	if err != nil {
		return err
	}
	containerID = strings.TrimSpace(containerID)
	if containerID == "" {
		return fmt.Errorf("mysql container was not created; check docker compose logs %s", opts.service)
	}

	fmt.Println("Waiting for MySQL healthcheck...")
	if err := waitForHealthy(containerID, opts.timeout); err != nil {
		return err
	}

	fmt.Println("Executing database/init.sql inside the MySQL container...")
	if err := run(projectRoot, "docker", "compose", "exec", "-T", opts.service, "sh", "-c", `mysql -uroot -p"$MYSQL_ROOT_PASSWORD" "$MYSQL_DATABASE" < /docker-entrypoint-initdb.d/01-init.sql`); err != nil {
		return err
	}

	fmt.Println("Database initialization completed.")
	return nil
}

func resolveProjectRoot(flagValue string) (string, error) {
	if flagValue != "" {
		return filepath.Abs(flagValue)
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if fileExists(filepath.Join(wd, "docker-compose.yml")) && fileExists(filepath.Join(wd, "database", "init.sql")) {
			return wd, nil
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			break
		}
		wd = parent
	}

	return "", errors.New("could not find project root containing docker-compose.yml and database/init.sql")
}

func waitForHealthy(containerID string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		status, err := output("", "docker", "inspect", "--format", "{{if .State.Health}}{{.State.Health.Status}}{{else}}{{.State.Status}}{{end}}", containerID)
		if err != nil {
			return err
		}
		status = strings.TrimSpace(status)
		if status == "healthy" {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("mysql did not become healthy before timeout; last status: %s", status)
		}
		time.Sleep(2 * time.Second)
	}
}

func run(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s %s failed: %w", name, strings.Join(args, " "), err)
	}
	return nil
}

func output(dir string, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%s %s failed: %w", name, strings.Join(args, " "), err)
	}
	return stdout.String(), nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
