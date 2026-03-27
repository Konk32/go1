package cmd

import (
    "fmt"
    "net/http"
    "os"
    "time"

    "github.com/spf13/cobra"
)

var timeout int

var checkCmd = &cobra.Command{
    Use:   "check <url>",
    Short: "Check if a URL responds",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        url := args[0]
        fmt.Printf("Checking %s (timeout: %ds)...\n", url, timeout)

        client := &http.Client{
            Timeout: time.Duration(timeout) * time.Second,
        }

        resp, err := client.Get(url)
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            os.Exit(1)
        }
        defer resp.Body.Close()

        fmt.Printf("Status: %d %s\n", resp.StatusCode, http.StatusText(resp.StatusCode))
    },
}

func init() {
    checkCmd.Flags().IntVarP(&timeout, "timeout", "t", 5, "Timeout in seconds")
    rootCmd.AddCommand(checkCmd)
}