package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ericm/stonks/api"
	"github.com/ericm/stonks/graph"
	"github.com/piquette/finance-go/datetime"
	"github.com/spf13/cobra"
)

func main() {
	var intervalCmd *string
	rootCmd := &cobra.Command{
		Use:   "stonks",
		Short: "A stock visualizer",
		Long:  "Displays realtime stocks in graph format in a terminal",
		Run: func(cmd *cobra.Command, args []string) {
			for _, symbol := range args {
				var interval datetime.Interval
				if intervalCmd == nil {
					interval = datetime.FifteenMins
				} else {
					interval = datetime.Interval(*intervalCmd)
				}
				chart, err := api.GetChart(strings.ToUpper(symbol), interval, nil)
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
				g, _ := graph.GenerateGraph(chart, 80, 12)
				fmt.Print(g)
			}
		},
	}
	intervalCmd = rootCmd.PersistentFlags().StringP("interval", "i", "15m", "stonks -t X[m|h]")
	rootCmd.Execute()
}
