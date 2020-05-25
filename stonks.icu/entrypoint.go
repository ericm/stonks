package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

func main() {
	configure()
	http.HandleFunc("/", handleSymbol)
	http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("port")), nil)
}

func handleSymbol(w http.ResponseWriter, r *http.Request) {
	symbols := strings.Split(r.URL.Path, "/")
	if r.Header.Get("User-Agent") != "" {
		w.Header().Add("Location", "https://github.com/ericm/stonks")
		w.WriteHeader(302)
		w.Write([]byte(" "))
	}
	for _, symbol := range symbols {
		fmt.Println(symbol)
	}
}
