package commands

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/homekit/homekit-cli/internal/assets"
)

// NewAssetsCommand provides subcommands for interacting with embedded assets.
func NewAssetsCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "assets",
		Short: "Inspect and manage embedded assets",
	}

	listCmd := &cobra.Command{
		Use:   "list [scripts|templates]",
		Args:  cobra.ExactArgs(1),
		Short: "List available assets by namespace",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listAssets(cmd, args[0])
		},
	}

	extractCmd := &cobra.Command{
		Use:   "extract [scripts|templates] <name> <dest>",
		Args:  cobra.ExactArgs(3),
		Short: "Extract an asset to a directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			return extractAsset(cmd, args[0], args[1], args[2])
		},
	}

	verifyCmd := &cobra.Command{
		Use:   "verify [scripts|templates] <name>",
		Args:  cobra.ExactArgs(2),
		Short: "Calculate a checksum for an asset",
		RunE: func(cmd *cobra.Command, args []string) error {
			return verifyAsset(cmd, args[0], args[1])
		},
	}

	root.AddCommand(listCmd, extractCmd, verifyCmd)
	return root
}

func listAssets(cmd *cobra.Command, namespace string) error {
	rt, err := runtimeFrom(cmd)
	if err != nil {
		return err
	}
	manager := assets.NewManager(assets.Embedded(), overrideDirectory(rt.Config))
	names, err := manager.List(namespace)
	if err != nil {
		return err
	}
	for _, name := range names {
		fmt.Fprintln(cmd.OutOrStdout(), name)
	}
	return nil
}

func extractAsset(cmd *cobra.Command, namespace, name, dest string) error {
	rt, err := runtimeFrom(cmd)
	if err != nil {
		return err
	}

	manager := assets.NewManager(assets.Embedded(), overrideDirectory(rt.Config))
	path, err := manager.Export(namespace, name, dest)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "exported %s to %s\n", name, filepath.ToSlash(path))
	return nil
}

func verifyAsset(cmd *cobra.Command, namespace, name string) error {
	rt, err := runtimeFrom(cmd)
	if err != nil {
		return err
	}

	manager := assets.NewManager(assets.Embedded(), overrideDirectory(rt.Config))
	sum, err := manager.Verify(namespace, name)
	if err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "%s\n", sum)
	return nil
}
