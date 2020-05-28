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

var (
	interval,
	save,
	remove,
	name *string
	year    *bool
	ytd     *bool
	week    *bool
	version *bool
	extra   *bool
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

	// Config defaults
	viper.SetDefault("config.standalone_height", 12)
	viper.SetDefault("config.favourites_height", 12)
	viper.SetDefault("config.default_theme", graph.LineTheme)

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
				fmt.Println(api.Version)
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
				if _, err := api.GetChart(saveCmd, datetime.FifteenMins, nil, nil, false); err != nil {
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

			chartTheme := graph.ChartTheme(viper.GetInt("config.default_theme"))

			switch *theme {
			case "line":
				chartTheme = graph.LineTheme
			case "dot":
				chartTheme = graph.DotTheme
			case "icon":
				chartTheme = graph.IconTheme
			default:
				if len(*theme) > 0 {
					fmt.Println("Unknown theme, must be \"line\", \"dot\" or \"icon\"")
					os.Exit(1)
				}
			}

			if len(args) == 0 {
				intervalCmd, start, end := parseTimeRange()
				extraCmd := false
				if end == nil {
					extraCmd = *extra
				}
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
					chart, err := api.GetChart(strings.ReplaceAll(strings.ToUpper(symbol), "_", "."), intervalCmd, start, end, extraCmd)
					if err != nil {
						fmt.Println(err.Error())
						continue
					}
					g, _ := graph.GenerateGraph(chart, 80, viper.GetInt("config.favourites_height"), chartTheme)
					fmt.Print(g)
				}
			}

			for _, symbol := range args {
				intervalCmd, start, end := parseTimeRange()
				extraCmd := false
				if end == nil {
					extraCmd = *extra
				}
				chart, err := api.GetChart(strings.ToUpper(symbol), intervalCmd, start, end, extraCmd)
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
				g, _ := graph.GenerateGraph(chart, 80, viper.GetInt("config.standalone_height"), chartTheme)
				fmt.Print(g)
			}
		},
	}
	interval = rootCmd.PersistentFlags().StringP("interval", "i", "15m", "stonks -i X[m|h] (eg 15m, 5m, 1h, 1d)")
	year = rootCmd.PersistentFlags().BoolP("year", "y", false, "Display the last year (will set interval to 5d)")
	ytd = rootCmd.PersistentFlags().Bool("ytd", false, "Display the year to date (will set interval to 5d)")
	week = rootCmd.PersistentFlags().BoolP("week", "w", false, "Display the last week (will set interval to 1d)")
	days = rootCmd.PersistentFlags().IntP("days", "d", 0, "24 hour period of stocks from X of days ago.")
	theme = rootCmd.PersistentFlags().StringP("theme", "t", "", "Display theme for the chart (Options: \"line\", \"dot\", \"icon\")")
	save = rootCmd.PersistentFlags().StringP("save", "s", "", "Add an item to the default stonks command. (Eg: -s AMD -n \"Advanced Micro Devices\")")
	remove = rootCmd.PersistentFlags().StringP("remove", "r", "", "Remove an item from favourites")
	name = rootCmd.PersistentFlags().StringP("name", "n", "", "Optional name for a stonk save")
	version = rootCmd.PersistentFlags().BoolP("version", "v", false, "stonks version")
	extra = rootCmd.PersistentFlags().BoolP("extra", "e", false, "Include extra pre + post time. (Only works for day)")

	rootCmd.Execute()
}

func parseTimeRange() (datetime.Interval, *datetime.Datetime, *datetime.Datetime) {
	var (
		intervalCmd datetime.Interval
		start       *datetime.Datetime
		end         *datetime.Datetime
	)
	switch {
	case *year:
		intervalCmd = datetime.FiveDay
		rn := time.Now()
		e := rn.AddDate(-1, 0, 0)
		start = datetime.New(&e)
		end = datetime.New(&rn)
	case *ytd:
		intervalCmd = datetime.FiveDay
		rn := time.Now()
		e := rn.AddDate(0, -int(rn.Month()), -rn.Day())
		start = datetime.New(&e)
		end = datetime.New(&rn)
	case *week:
		intervalCmd = datetime.OneHour
		rn := time.Now()
		e := rn.AddDate(0, 0, -7)
		start = datetime.New(&e)
		end = datetime.New(&rn)
	case interval == nil:
		intervalCmd = datetime.FifteenMins
	default:
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
	return intervalCmd, start, end
}
