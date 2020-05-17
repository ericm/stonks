package graph

import (
	"fmt"

	"github.com/ericm/stonks/api"
	"github.com/shopspring/decimal"
)

func borderHorizontal(out *string, width int) {
	for _i := 0; _i < width-2; _i++ {
		*out += "╌"
	}
}

// GenerateGraph with ASCII graph with ANSI escapes
func GenerateGraph(chart *api.Chart, width int, height int) (string, error) {
	out := "┌"
	borderHorizontal(&out, width)
	out += "┐"
	// interval, err := time.ParseDuration(string(chart.Interval))
	// if err != nil {
	// 	return "", err
	// }
	matrix := make([][]*api.Bar, height)
	for i := range matrix {
		matrix[i] = make([]*api.Bar, width)
	}
	ran := chart.High.Sub(chart.Low)
	for _, bar := range chart.Bars {
		fmt.Println(bar.Current.Sub(chart.Low).Div(ran).Mul(decimal.NewFromInt((int64(height)))).Floor().IntPart())
	}
	return out, nil
}
