package cmd

import (
    "fmt"
    "net/http"
    "time"
    "sync"

    "github.com/spf13/cobra"
)

var timeout int

type Result struct {
    URL        string
    StatusCode int
    StatusText string
    Elapsed    time.Duration
    Err        error
}

var checkCmd = &cobra.Command{
    Use:   "check <url> [url...]",
    Short: "Check if one or more URLs respond",
    Args:  cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        client := &http.Client{
            Timeout: time.Duration(timeout) * time.Second,
        }

        results := make(chan Result, len(args))
        var wg sync.WaitGroup

        for _, url:= range args {
            wg.Add(1)
            go func(u string) {
                defer wg.Done()
                results <- checkURL(client, u)
            }(url)
        }

        // Close the channel once all goroutines finish
        go func() {
            wg.Wait()
            close(results)
        }()

        for result := range results {
            if result.Err != nil {
                fmt.Printf("%-35s ERROR — %v\n", result.URL, result.Err)
                continue
            }
            fmt.Printf("%-35s %d %s — %dms\n",
                result.URL,
                result.StatusCode,
                result.StatusText,
                result.Elapsed.Milliseconds(),
            )
        }
    },
}
func checkURL(client *http.Client, url string) Result {
    start := time.Now()

    resp, err := client.Get(url)
    if err != nil {
        return Result{URL: url, Err: err}
    }
    defer resp.Body.Close()

    return Result{
        URL:        url,
        StatusCode: resp.StatusCode,
        StatusText: http.StatusText(resp.StatusCode),
        Elapsed:    time.Since(start),
    }
}

func init() {
    checkCmd.Flags().IntVarP(&timeout, "timeout", "t", 5, "Timeout in seconds")
    rootCmd.AddCommand(checkCmd)
}