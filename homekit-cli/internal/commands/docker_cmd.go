package commands

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	executor "github.com/homekit/homekit-cli/internal/exec"
)

// NewDockerCommand provides helper commands around Docker tooling.
func NewDockerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "docker",
		Short: "Quality-of-life helpers for Docker workflows",
	}

	cmd.AddCommand(newDockerPruneCommand())
	cmd.AddCommand(newDockerImagesUpdateCommand())
	return cmd
}

func newDockerPruneCommand() *cobra.Command {
	var aggressive bool

	c := &cobra.Command{
		Use:   "prune",
		Short: "Prune unused Docker resources safely",
		RunE: func(cmd *cobra.Command, args []string) error {
			rt, err := runtimeFrom(cmd)
			if err != nil {
				return err
			}

			flags := []string{"system", "prune"}
			if aggressive {
				flags = append(flags, "--all")
			}
			flags = append(flags, "--force")

			spec := executor.Spec{
				Command:       "docker",
				Args:          flags,
				Timeout:       15 * time.Minute,
				CaptureOutput: !rt.DryRun,
				DryRun:        rt.DryRun,
			}

			rt.Logger.Info().Strs("args", spec.Args).Msg("running docker prune")

			_, err = executor.Run(cmd.Context(), spec)
			return err
		},
	}

	c.Flags().BoolVar(&aggressive, "aggressive", false, "Include unused images in prune")
	return c
}

func newDockerImagesUpdateCommand() *cobra.Command {
	var pullAll bool

	c := &cobra.Command{
		Use:   "images update [image...]",
		Short: "Pull updated images for provided repositories",
		RunE: func(cmd *cobra.Command, args []string) error {
			rt, err := runtimeFrom(cmd)
			if err != nil {
				return err
			}
			images := args
			if len(images) == 0 && !pullAll {
				return fmt.Errorf("specify image names or use --all")
			}

			if pullAll {
				images = append(images, "hello-world") // placeholder
			}

			for _, image := range images {
				rt.Logger.Info().Str("image", image).Msg("pulling docker image")
				_, err := executor.Run(cmd.Context(), executor.Spec{
					Command:       "docker",
					Args:          []string{"pull", image},
					Timeout:       10 * time.Minute,
					CaptureOutput: !rt.DryRun,
					DryRun:        rt.DryRun,
				})
				if err != nil {
					return err
				}
			}
			return nil
		},
	}

	c.Flags().BoolVar(&pullAll, "all", false, "Pull all configured images (placeholder)")
	return c
}
