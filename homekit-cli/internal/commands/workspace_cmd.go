package commands

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/homekit/homekit-cli/internal/assets"
	"github.com/homekit/homekit-cli/internal/core"
	"github.com/homekit/homekit-cli/internal/util/pathformat"
	"github.com/homekit/homekit-cli/internal/util/templateutil"
	"github.com/spf13/cobra"
)

var (
	imageTypes = []string{"go1.23.2", "uv", "nvm", "default"}
)

type WorkspaceOptions struct {
	DirPath string `mapstructure:"dir_path"`
	Name    string `mapstructure:"name"`
	Type    string `mapstructure:"type"`
}

func NewWorkspaceCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "workspace [new]",
		Short: "Manage workspaces",
	}

	newCmd := &cobra.Command{
		Use:   "new -d <directory> -n <name> -t <type>",
		Short: "Create a new workspace with the given language base",
		RunE: func(cmd *cobra.Command, args []string) error {
			rt, err := runtimeFrom(cmd)
			if err != nil {
				return err
			}

			imageType, _ := cmd.Flags().GetString("t")
			if imageType == "" {
				imageType = "default"
			}
			dirStr, _ := cmd.Flags().GetString("d")
			name, _ := cmd.Flags().GetString("n")

			if !slices.Contains(imageTypes, imageType) {
				return fmt.Errorf("invalid image type: %s. Please choose from %s", imageType, strings.Join(imageTypes, ", "))
			}

			opts := WorkspaceOptions{
				DirPath: pathformat.RenderFullPath(dirStr),
				Name:    name,
				Type:    imageType,
			}

			workspaceDir, err := createWorkspaceSkeleton(rt, opts)
			if err != nil {
				return err
			}

			rt.Logger.Info().Msgf("Workspace created successfully in %s", workspaceDir)
			return nil
		},
		Args: cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
	}

	newCmd.Flags().StringP("dir", "d", ".", "Directory to create the workspace in")
	newCmd.Flags().StringP("name", "n", "", "Name of the workspace")
	newCmd.Flags().StringP("type", "t", "default", "Type of the workspace")

	root.AddCommand(newCmd)
	return root
}

/*
*
Create a new workspace skeleton.

The skeleton will be created in the given directory.
It will contain the following files:
- (dir if given, otherwise current directory)
  - README.md
  - code (dir)
  - Makefile
*/
func createWorkspaceSkeleton(rt *core.Runtime, opts WorkspaceOptions) (string, error) {
	// init workspace dir
	workspaceDir := ""
	if opts.DirPath == "" {
		workspaceDir = pathformat.Pwd()
	} else {
		workspaceDir = pathformat.RenderFullPath(opts.DirPath)
		err := pathformat.MakeDirIfNotExists(workspaceDir)
		if err != nil {
			return "", err
		}
	}
	if opts.Name == "" {
		opts.Name = pathformat.Base(workspaceDir)
	}

	rt.Logger.Info().Msgf("Name: %s", opts.Name)
	rt.Logger.Info().Msgf("Type: %s", opts.Type)

	// create README.md with template replacement
	assetManager := assets.NewManager(assets.Embedded(), "")
	readmeContent, err := assetManager.OpenBytes(assets.AssetNamespaceWorkspaces.String(), "README.md")
	if err != nil {
		return "", err
	}

	readmePath := pathformat.Join(workspaceDir, "README.md")
	readmeFinalContent, err := templateutil.RenderTemplateInBytes(readmeContent, opts, "readme", rt.BufPool)
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(readmePath, readmeFinalContent, 0644); err != nil {
		return "", err
	}

	rt.Logger.Info().Msgf("README.md created in %s", readmePath)

	// create code directory
	codeDir := pathformat.Join(workspaceDir, "code")
	if err := os.MkdirAll(codeDir, 0755); err != nil {
		return "", err
	}

	rt.Logger.Info().Msgf("Code directory created in %s", codeDir)

	// create Makefile
	makefileContent, err := assetManager.OpenBytes(assets.AssetNamespaceWorkspaces.String(), "Makefile")
	if err != nil {
		return "", err
	}
	makefilePath := pathformat.Join(workspaceDir, "Makefile")
	makefileFinalContent, err := templateutil.RenderTemplateInBytes(makefileContent, opts, "Makefile", rt.BufPool)
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(makefilePath, makefileFinalContent, 0644); err != nil {
		return "", err
	}

	rt.Logger.Info().Msgf("Makefile created in %s", makefilePath)

	// create compose.dev.yml
	composeDevContent, err := assetManager.OpenBytes(assets.AssetNamespaceWorkspaces.String(), "compose.dev.yml")
	if err != nil {
		return "", err
	}
	composeDevFinalContent, err := templateutil.RenderTemplateInBytes(composeDevContent, opts, "compose.dev.yml", rt.BufPool)
	if err != nil {
		return "", err
	}
	composeDevPath := pathformat.Join(workspaceDir, "compose.dev.yml")
	if err := os.WriteFile(composeDevPath, composeDevFinalContent, 0644); err != nil {
		return "", err
	}
	rt.Logger.Info().Msgf("compose.dev.yml created")

	return workspaceDir, nil
}
