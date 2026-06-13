package component

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/goccy/go-yaml"
	"github.com/slingshot-pipelines/cli/internal/components"
	"github.com/spf13/cobra"
)

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List components in the repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			h := &listHarness{}

			componentsPath, err := h.findComponentsDir()
			if err != nil {
				return err
			}

			componentFiles, err := h.listComponentFiles(componentsPath)
			if err != nil {
				return err
			}

			parsedComponents, err := h.parseComponents(componentFiles)
			if err != nil {
				return err
			}

			out, err := json.MarshalIndent(parsedComponents, "", "  ")

			fmt.Println(string(out))

			return nil
		},
	}

	return cmd
}

type listHarness struct {
}

func (h *listHarness) findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error determining working directory: %w", err)
	}

	cursor := dir
	for path.Dir(cursor) != cursor {
		_, err := os.Stat(path.Join(cursor, ".git"))
		if err == nil {
			return cursor, nil
		} else if !errors.Is(err, os.ErrNotExist) {
			return "", err
		}

		cursor = path.Dir(cursor)
	}

	return "", fmt.Errorf("no .git directory found in any parent directory")
}

func (h *listHarness) findComponentsDir() (string, error) {
	repoRoot, err := h.findRepoRoot()
	if err != nil {
		return "", fmt.Errorf("error determining repository root: %w", err)
	}

	componentsPath := path.Join(repoRoot, ".components")
	stats, err := os.Stat(componentsPath)
	if err != nil {
		return "", err
	}

	if !stats.IsDir() {
		return "", fmt.Errorf("expected a directory at %s, got a file", componentsPath)
	}

	return componentsPath, nil
}

func (h *listHarness) listComponentFiles(componentsDir string) ([]string, error) {
	componentFiles := []string{}

	entries, err := os.ReadDir(componentsDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			infoFile, err := h.findComponentFile(path.Join(componentsDir, entry.Name()))
			if err != nil {
				return nil, err
			}

			if infoFile != "" {
				componentFiles = append(componentFiles, infoFile)
			}
		}
	}

	return componentFiles, nil
}

func (h *listHarness) findComponentFile(componentDir string) (string, error) {
	entries, err := os.ReadDir(componentDir)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if entry.Name() == "info.yml" || entry.Name() == "info.yaml" {
			return path.Join(componentDir, entry.Name()), nil
		}
	}

	return "", nil
}

func (h *listHarness) parseComponents(componentFiles []string) ([]components.Component, error) {
	parsedComponents := []components.Component{}

	for _, file := range componentFiles {
		parsedComponent, err := h.parseComponentInfo(file)
		if err != nil {
			return nil, err
		}

		parsedComponents = append(parsedComponents, *parsedComponent)
	}

	return parsedComponents, nil
}

func (h *listHarness) parseComponentInfo(infoFile string) (*components.Component, error) {
	contents, err := os.ReadFile(infoFile)
	if err != nil {
		return nil, err
	}

	component := components.Component{}
	err = yaml.Unmarshal(contents, &component)
	if err != nil {
		return nil, err
	}

	return &component, nil
}
