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
	rootCmd := &cobra.Command{
		Use:   "stonks",
		Short: "A stock visualizer",
		Long:  "Displays realtime stocks in graph format in a terminal",
		Run: func(cmd *cobra.Command, args []string) {
			for _, symbol := range args {
				chart, err := api.GetChart(strings.ToUpper(symbol), datetime.FifteenMins, nil)
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
				g, _ := graph.GenerateGraph(chart, 80, 12)
				fmt.Print(g)
			}
		},
	}
	rootCmd.SetUsageTemplate("Usage:\n  stonks [flags] [symbols]\n")
	rootCmd.Execute()
}
