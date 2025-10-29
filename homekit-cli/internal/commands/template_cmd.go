package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/homekit/homekit-cli/internal/assets"
	"github.com/homekit/homekit-cli/internal/templating"
)

// NewTemplateCommand exposes templating workflows.
func NewTemplateCommand() *cobra.Command {
	var dataFiles []string
	var output string

	cmd := &cobra.Command{
		Use:   "template",
		Short: "Render templates from embedded assets or local files",
	}

	renderCmd := &cobra.Command{
		Use:   "render <template>",
		Args:  cobra.ExactArgs(1),
		Short: "Render an embedded template to stdout or a file",
		RunE: func(cmd *cobra.Command, args []string) error {
			rt, err := runtimeFrom(cmd)
			if err != nil {
				return err
			}

			manager := assets.NewManager(assets.Embedded(), overrideDirectory(rt.Config))

			data, err := aggregateTemplateData(dataFiles)
			if err != nil {
				return err
			}

			renderer := templating.Renderer{}
			var out *os.File
			if output != "" {
				dest, err := os.Create(output)
				if err != nil {
					return err
				}
				defer dest.Close()
				out = dest
			} else {
				out = os.Stdout
			}

			handle, err := manager.Open("templates", args[0])
			if err != nil {
				return err
			}
			defer handle.Close()

			if err := renderer.Render(handle, data, out); err != nil {
				return err
			}
			return nil
		},
	}

	renderCmd.Flags().StringSliceVarP(&dataFiles, "data", "d", nil, "YAML data files to merge")
	renderCmd.Flags().StringVarP(&output, "output", "o", "", "Destination file (default stdout)")

	cmd.AddCommand(renderCmd)
	return cmd
}

func aggregateTemplateData(paths []string) (map[string]any, error) {
	out := map[string]any{}
	for _, path := range paths {
		file, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read data file %s: %w", path, err)
		}
		doc := map[string]any{}
		if err := yaml.Unmarshal(file, &doc); err != nil {
			return nil, fmt.Errorf("parse data file %s: %w", path, err)
		}
		for k, v := range doc {
			out[k] = v
		}
	}
	return out, nil
}
