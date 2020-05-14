package main

import (
	"fmt"

	"github.com/ericm/stonks/api"
	"github.com/piquette/finance-go/datetime"
)

func main() {
	chart, _ := api.GetChart("AMD", datetime.FifteenMins, nil)
	fmt.Println(chart)
}
