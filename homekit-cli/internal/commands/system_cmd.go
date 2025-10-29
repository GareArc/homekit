package commands

import (
	"fmt"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/spf13/cobra"
)

// NewSystemCommand delivers lightweight health checks.
func NewSystemCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sys",
		Short: "Inspect local system health metrics",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "health",
		Short: "Display basic system health information",
		RunE: func(cmd *cobra.Command, args []string) error {
			memInfo, err := mem.VirtualMemory()
			if err != nil {
				return err
			}

			loadAvg, err := load.Avg()
			if err != nil {
				return err
			}

			diskInfo, err := disk.Usage("/")
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Time: %s\n", time.Now().Format(time.RFC3339))
			fmt.Fprintf(cmd.OutOrStdout(), "OS: %s / %s\n", runtime.GOOS, runtime.GOARCH)
			fmt.Fprintf(cmd.OutOrStdout(), "Memory: %.2f%% used (%.2f GiB / %.2f GiB)\n",
				memInfo.UsedPercent,
				float64(memInfo.Used)/1_073_741_824,
				float64(memInfo.Total)/1_073_741_824,
			)
			fmt.Fprintf(cmd.OutOrStdout(), "Load Average: %.2f %.2f %.2f\n", loadAvg.Load1, loadAvg.Load5, loadAvg.Load15)
			fmt.Fprintf(cmd.OutOrStdout(), "Disk: %.2f%% used (%.2f GiB / %.2f GiB)\n",
				diskInfo.UsedPercent,
				float64(diskInfo.Used)/1_073_741_824,
				float64(diskInfo.Total)/1_073_741_824,
			)
			return nil
		},
	})

	return cmd
}
