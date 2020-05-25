package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ericm/stonks/api"
	"github.com/ericm/stonks/graph"
	"github.com/piquette/finance-go/datetime"
	"github.com/spf13/viper"
)

const version = "1.0.6"
const footer = "\nLike Stonks? Star it on GitHub: https://github.com/ericm/stonks\nstonks " + version + "\n"

func main() {
	configure()
	http.HandleFunc("/", handleSymbol)
	http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("port")), nil)
}

func handleSymbol(w http.ResponseWriter, r *http.Request) {
	symbols := strings.Split(r.URL.Path, "/")
	plainText := false
	for _, client := range clients {
		if strings.Contains(r.Header.Get("User-Agent"), client+"/") {
			plainText = true
			break
		}
	}
	if !plainText {
		w.Header().Add("Location", "https://github.com/ericm/stonks")
		w.WriteHeader(302)
		w.Write([]byte(" "))
		return
	}
	output := ""
	num := 0
	for _, symbol := range symbols {
		if len(symbol) > 0 {
			num++
			symbol = strings.ToUpper(symbol)
			chart, err := api.GetChart(symbol, datetime.FifteenMins, nil, nil, false)
			if err != nil {
				w.WriteHeader(403)
				w.Write([]byte(err.Error() + "\n" + footer))
				return
			}
			out, err := graph.GenerateGraph(chart, 80, 12, graph.LineTheme)
			if err != nil {
				w.WriteHeader(403)
				w.Write([]byte(err.Error() + "\n" + footer))
				return
			}
			output += out
		}
	}
	w.WriteHeader(200)
	if num == 0 {
		w.Write([]byte("Please provide stonks in the format:\nstonks.icu/amd/intl\n" + footer))
	} else {
		w.Write([]byte(output + footer))
	}
}
