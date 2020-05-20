package main

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/ericm/stonks/api"
	"github.com/ericm/stonks/graph"
	"github.com/piquette/finance-go/datetime"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const Version = "1.0.1"

var (
	interval,
	save,
	remove,
	name *string
	week    *bool
	version *bool
	theme   *string
	days    *int

	configPath string
)

func setDefaults() {
	home, ok := os.LookupEnv("HOME")
	if !ok {
		panic("No HOME env var set")
	}
	configPath = fmt.Sprintf("%s/.config", home)
	viper.AddConfigPath(configPath)
	viper.SetConfigName("stonks.yml")
	viper.SetConfigType("yaml")

	viper.SetDefault("favourites", map[string]interface{}{})

	viper.ReadInConfig()
}

func main() {
	setDefaults()
	rootCmd := &cobra.Command{
		Use:   "stonks",
		Short: "A stock visualizer",
		Long:  "Displays realtime stocks in graph format in a terminal",
		Run: func(cmd *cobra.Command, args []string) {
			if *version {
				fmt.Println(Version)
				return
			}

			if len(*remove) > 0 {
				saveCmd := strings.ToLower(*remove)
				favourites, ok := viper.Get("favourites").(map[string]interface{})
				if !ok {
					fmt.Println("Read config error")
					os.Exit(1)
				}
				delete(favourites, saveCmd)
				viper.Set("favourites", favourites)
				viper.WriteConfig()
				return
			}

			if len(*save) > 0 {
				saveCmd := strings.ToUpper(*save)
				if _, err := api.GetChart(saveCmd, datetime.FifteenMins, nil, nil); err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
				favourites, ok := viper.Get("favourites").(map[string]interface{})
				if !ok {
					fmt.Println("Read config error")
					os.Exit(1)
				}
				nameCmd := saveCmd
				if len(*name) > 0 {
					nameCmd = *name
				}
				favourites[strings.ReplaceAll(saveCmd, ".", "_")] = nameCmd
				viper.Set("favourites", favourites)
				if err := viper.WriteConfig(); err != nil {
					err = viper.WriteConfigAs(path.Join(configPath, "stonks.yml"))
					if err != nil {
						fmt.Println(err.Error())
						os.Exit(1)
					}
				}
				return
			}

			chartTheme := graph.LineTheme

			switch {
			case *theme == "line":
				chartTheme = graph.LineTheme
			case *theme == "dot":
				chartTheme = graph.DotTheme
			case *theme == "icon":
				chartTheme = graph.IconTheme
			default:
				fmt.Println("Unknown theme, must be \"line\", \"dot\" or \"icon\"")
				os.Exit(1)
			}

			if len(args) == 0 {
				// Favourites
				favourites, ok := viper.Get("favourites").(map[string]interface{})
				if !ok {
					fmt.Println("Read config error")
					os.Exit(1)
				}
				if len(favourites) == 0 {
					fmt.Println("No favourites added. You can add them in the format 'stonks -s AMD -n \"Advanced Micro Devices\"'")
				}

				keys := make([]string, 0, len(favourites))
				for k := range favourites {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, symbol := range keys {
					name := favourites[symbol].(string)
					fmt.Println(name + ":")
					chart, err := api.GetChart(strings.ReplaceAll(strings.ToUpper(symbol), "_", "."), datetime.FifteenMins, nil, nil)
					if err != nil {
						fmt.Println(err.Error())
						continue
					}
					g, _ := graph.GenerateGraph(chart, 80, 12, chartTheme)
					fmt.Print(g)
				}
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
				g, _ := graph.GenerateGraph(chart, 80, 12, chartTheme)
				fmt.Print(g)
			}
		},
	}
	interval = rootCmd.PersistentFlags().StringP("interval", "i", "15m", "stonks -i X[m|h] (eg 15m, 5m, 1h, 1d)")
	week = rootCmd.PersistentFlags().BoolP("week", "w", false, "Display the last week (will set interval to 1d)")
	days = rootCmd.PersistentFlags().IntP("days", "d", 0, "Stocks from X number of days ago.")
	theme = rootCmd.PersistentFlags().StringP("theme", "t", "line", "Display theme for the chart (Options: \"line\", \"dot\", \"icon\")")
	save = rootCmd.PersistentFlags().StringP("save", "s", "", "Add an item to the default stonks command. (Eg: -s AMD -n \"Advanced Micro Devices\")")
	remove = rootCmd.PersistentFlags().StringP("remove", "r", "", "Remove an item from favourites")
	name = rootCmd.PersistentFlags().StringP("name", "n", "", "Optional name for a stonk save")
	version = rootCmd.PersistentFlags().BoolP("version", "v", false, "stonks version")

	rootCmd.Execute()
}
