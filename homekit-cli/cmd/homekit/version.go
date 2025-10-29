package homekit

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func newVersionCommand() *cobra.Command {
	var output string

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			info := map[string]string{
				"version": version,
				"commit":  commit,
				"date":    date,
				"source":  source,
			}

			switch output {
			case "json":
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(info)
			default:
				tw := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
				for k, v := range info {
					fmt.Fprintf(tw, "%s:\t%s\n", k, v)
				}
				return tw.Flush()
			}
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "table", "Output format (table|json)")
	return cmd
}
