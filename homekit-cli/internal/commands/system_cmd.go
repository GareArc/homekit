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
			rt, err := runtimeFrom(cmd)
			if err != nil {
				return err
			}

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

			rt.Logger.Info().
				Time("time", time.Now()).
				Str("os", fmt.Sprintf("%s / %s", runtime.GOOS, runtime.GOARCH)).
				Str("memory", fmt.Sprintf("%.2f%% used (%.2f GiB / %.2f GiB)", memInfo.UsedPercent, float64(memInfo.Used)/1_073_741_824, float64(memInfo.Total)/1_073_741_824)).
				Str("load_average", fmt.Sprintf("%.2f %.2f %.2f", loadAvg.Load1, loadAvg.Load5, loadAvg.Load15)).
				Str("disk", fmt.Sprintf("%.2f%% used (%.2f GiB / %.2f GiB)", diskInfo.UsedPercent, float64(diskInfo.Used)/1_073_741_824, float64(diskInfo.Total)/1_073_741_824)).
				Msg("System health")
			return nil
		},
	})

	return cmd
}
