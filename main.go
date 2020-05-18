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
	"github.com/spf13/viper"
)

var (
	interval *string
	week     *bool
	days     *int
)

func main() {
	home, ok := os.LookupEnv("HOME")
	if !ok {
		panic("No HOME env var set")
	}
	viper.AddConfigPath(fmt.Sprintf("%s/.config/stonks/", home))
	viper.SetConfigName("favourites")
	viper.SetConfigType("yaml")
	viper.ReadInConfig()
	rootCmd := &cobra.Command{
		Use:   "stonks",
		Short: "A stock visualizer",
		Long:  "Displays realtime stocks in graph format in a terminal",
		Run: func(cmd *cobra.Command, args []string) {
			for _, symbol := range args {
				var intervalCmd datetime.Interval
				var start *datetime.Datetime
				var end *datetime.Datetime
				if *week {
					intervalCmd = datetime.OneHour
					rn := time.Now()
					e := rn.AddDate(0, 0, -7)
					start = datetime.New(&e)
					end = datetime.New(&rn)
				} else if interval == nil {
					intervalCmd = datetime.FifteenMins
				} else {
					intervalCmd = datetime.Interval(*interval)
				}
				if *days > 0 {
					s := time.Now().AddDate(0, 0, *days*-1)
					y, m, d := s.Date()
					s = time.Date(y, m, d, 0, 0, 0, 0, s.Location())

					start = datetime.New(&s)
					e := time.Date(y, m, d, 23, 0, 0, 0, s.Location())
					end = datetime.New(&e)
				}
				chart, err := api.GetChart(strings.ToUpper(symbol), intervalCmd, start, end)
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
				g, _ := graph.GenerateGraph(chart, 80, 12)
				fmt.Print(g)
			}
		},
	}
	interval = rootCmd.PersistentFlags().StringP("interval", "i", "15m", "stonks -t X[m|h] (eg 15m, 5m, 1h, 1d)")
	week = rootCmd.PersistentFlags().BoolP("week", "w", false, "Display the last week (will set interval to 1d)")
	days = rootCmd.PersistentFlags().IntP("days", "d", 0, "Stocks from X number of days ago.")
	rootCmd.Execute()
}
