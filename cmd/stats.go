package cmd

import (
	"fmt"
	"math"

	"github.com/spf13/cobra"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

var statsCmd = &cobra.Command{
	Use: "stats",
	Short: "Show current system resource usage",
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		printMemory()
		printCPU()
		printDisk()
	},
}

func printMemory() {
	v, err := mem.VirtualMemory()
	if err != nil {
		fmt.Printf("Memory error: %v\n", err)
		return
	}

	fmt.Println("── Memory ───────────────────────────")
	fmt.Printf("  Total:     %s\n", formatBytes(v.Total))
	fmt.Printf("  Used:      %s (%.1f%%)\n", formatBytes(v.Used), v.UsedPercent)
	fmt.Printf("  Available: %s\n", formatBytes(v.Available))
}

func printCPU() {
	percentages, err := cpu.Percent(0, false)
	if err != nil {
		fmt.Printf("CPU error: %v\n", err)
		return
	}

	fmt.Println("── CPU ──────────────────────────────")
	fmt.Printf("  Usage: %.1f%%\n", percentages[0])
}

func printDisk() {
	usage, err := disk.Usage("/")
	if err != nil {
		// Try Windows root if / fails
		usage, err = disk.Usage("C:")
		if err != nil {
			fmt.Printf("Disk error: %v\n", err)
			return
		}
	}

	fmt.Println("── Disk ─────────────────────────────")
	fmt.Printf("  Total:     %s\n", formatBytes(usage.Total))
	fmt.Printf("  Used:      %s (%.1f%%)\n", formatBytes(usage.Used), usage.UsedPercent)
	fmt.Printf("  Free:      %s\n", formatBytes(usage.Free))
}

func formatBytes(bytes uint64) string {
	units := []string{"B", "KB", "MB", "GB", "TB"}
	if bytes == 0 {
		return "0 B"
	}
	exp := int(math.Log(float64(bytes))/math.Log(1024))
	if exp >= len(units) {
		exp = len(units) - 1
	}
	value := float64(bytes) / math.Pow(1024, float64(exp))
	return fmt.Sprintf("%.1f %s", value, units[exp])
}

func init() {
	rootCmd.AddCommand(statsCmd)
}