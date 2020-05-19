package main

import (
	"fmt"
	"os"
	"path"
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
	save     *string
	name     *string
	week     *bool
	days     *int

	configPath string
)

func setDefaults() {
	home, ok := os.LookupEnv("HOME")
	if !ok {
		panic("No HOME env var set")
	}
	configPath = fmt.Sprintf("%s/.config", home)
	viper.AddConfigPath(configPath)
	viper.SetConfigName("stonks")
	viper.SetConfigType("yaml")

	viper.SetDefault("favourites", map[string]string{})

	viper.ReadInConfig()
}

func main() {
	setDefaults()
	rootCmd := &cobra.Command{
		Use:   "stonks",
		Short: "A stock visualizer",
		Long:  "Displays realtime stocks in graph format in a terminal",
		Run: func(cmd *cobra.Command, args []string) {
			if len(*save) > 0 {
				*save = strings.ToUpper(*save)
				if _, err := api.GetChart(strings.ToUpper(*save), datetime.FifteenMins, nil, nil); err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
				favourites, ok := viper.Get("favourites").(map[string]string)
				if !ok {
					fmt.Println("Read config error")
					os.Exit(1)
				}
				nameCmd := *save
				if len(*name) > 0 {
					nameCmd = *name
				}
				favourites[*save] = nameCmd
				viper.Set("favourites", favourites)
				if err := viper.WriteConfig(); err != nil {
					err = viper.WriteConfigAs(path.Join(configPath, "stonks"))
					if err != nil {
						fmt.Println(err.Error())
						os.Exit(1)
					}
				}
				return
			}
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
	save = rootCmd.PersistentFlags().StringP("save", "s", "", "Add an item to the default stonks command. (Eg: -s AMD -n \"Advanced Micro Devices\")")
	name = rootCmd.PersistentFlags().StringP("name", "n", "", "Optional name for a stonk save")

	rootCmd.Execute()
}
