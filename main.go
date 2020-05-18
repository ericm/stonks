package main

import (
	"fmt"
	"os"

	"github.com/ericm/stonks/api"
	"github.com/ericm/stonks/graph"
	"github.com/piquette/finance-go/datetime"
)

func main() {
	chart, err := api.GetChart("AMD", datetime.FifteenMins, nil)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	g, _ := graph.GenerateGraph(chart, 80, 10)
	fmt.Print(g)
}
