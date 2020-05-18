package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ericm/stonks/api"
	"github.com/ericm/stonks/graph"
	"github.com/piquette/finance-go/datetime"
	"github.com/spf13/cobra"
)

func main() {
	var (
		intervalCmd *string
		week        *bool
	)
	rootCmd := &cobra.Command{
		Use:   "stonks",
		Short: "A stock visualizer",
		Long:  "Displays realtime stocks in graph format in a terminal",
		Run: func(cmd *cobra.Command, args []string) {
			for _, symbol := range args {
				var interval datetime.Interval
				var start *datetime.Datetime
				var end *datetime.Datetime
				if *week {
					interval = datetime.OneHour
					rn := time.Now()
					e := rn.AddDate(0, 0, -7)
					start = datetime.New(&e)
					end = datetime.New(&rn)
				} else if intervalCmd == nil {
					interval = datetime.FifteenMins
				} else {
					interval = datetime.Interval(*intervalCmd)
				}
				chart, err := api.GetChart(strings.ToUpper(symbol), interval, start, end)
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
				g, _ := graph.GenerateGraph(chart, 80, 12)
				fmt.Print(g)
			}
		},
	}
	intervalCmd = rootCmd.PersistentFlags().StringP("interval", "i", "15m", "stonks -t X[m|h] (eg 15m, 5m, 1h, 1d)")
	week = rootCmd.PersistentFlags().BoolP("week", "w", false, "Display the last week (will set interval to 1d)")
	rootCmd.Execute()
}
