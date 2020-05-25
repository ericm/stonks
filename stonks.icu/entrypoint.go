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
	for _, symbol := range symbols {
		if len(symbol) > 0 {
			symbol = strings.ToUpper(symbol)
			chart, err := api.GetChart(symbol, datetime.FifteenMins, nil, nil, false)
			if err != nil {
				w.WriteHeader(403)
				w.Write([]byte(err.Error() + "\n"))
				return
			}
			out, err := graph.GenerateGraph(chart, 80, 12, graph.LineTheme)
			if err != nil {
				w.WriteHeader(403)
				w.Write([]byte(err.Error() + "\n"))
				return
			}
			output += out
		}
	}
	w.WriteHeader(200)
	w.Write([]byte(output))
}
