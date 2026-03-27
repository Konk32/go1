package cmd

import (
	"fmt"
	"net/http"
	"time"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
)

var interval int

var watchCmd = &cobra.Command{
	Use: "Watch [url]",
	Short: "Poll a URL repeadetly until stopped",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]

		// Channel that recieves a signal when the user hits Ctrl+C
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)

		client := &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		}

		fmt.Printf("Watching %s every %ds — press Ctrl+C to stop\n\n", url, interval)

		 ticker := time.NewTicker(time.Duration(interval) * time.Second)
        defer ticker.Stop()

        // Run a check immediately before the first tick
        doCheck(client, url)

        for {
            select {
            case <-ticker.C:
                doCheck(client, url)
            case <-quit:
                fmt.Println("\nStopped.")
                return
            }
        }
    },
}

func doCheck(client *http.Client, url string) {
    start := time.Now()

    resp, err := client.Get(url)
    if err != nil {
        fmt.Printf("[%s] ERROR — %v\n", time.Now().Format("15:04:05"), err)
        return
    }
    defer resp.Body.Close()

    elapsed := time.Since(start)
    fmt.Printf("[%s] %d %s — %dms\n",
        time.Now().Format("15:04:05"),
        resp.StatusCode,
        http.StatusText(resp.StatusCode),
        elapsed.Milliseconds(),
    )
}

func init() {
    watchCmd.Flags().IntVarP(&interval, "interval", "i", 5, "Seconds between checks")
    watchCmd.Flags().IntVarP(&timeout, "timeout", "t", 5, "Timeout in seconds")
    rootCmd.AddCommand(watchCmd)
}