package graph

import (
	"github.com/ericm/stonks/api"
	"github.com/shopspring/decimal"
)

func borderHorizontal(out *string, width int) {
	for _i := 0; _i < width-2; _i++ {
		*out += "━"
	}
}

// GenerateGraph with ASCII graph with ANSI escapes
func GenerateGraph(chart *api.Chart, width int, height int) (string, error) {
	out := "┏"
	borderHorizontal(&out, width)
	out += "┓"
	// interval, err := time.ParseDuration(string(chart.Interval))
	// if err != nil {
	// 	return "", err
	// }
	matrix := make([][]*api.Bar, height)
	for i := range matrix {
		matrix[i] = make([]*api.Bar, width)
	}
	ran := chart.High.Sub(chart.Low)
	spacing := (width) / (chart.Length)
	out += "\n"
	var last *api.Bar
	for x, bar := range chart.Bars {
		bar.Char = "─"
		y := int(bar.Current.Sub(chart.Low).Div(ran).Mul(
			decimal.NewFromInt((int64(height)))).Floor().IntPart())
		matrix[y][x*spacing] = bar
		bar.Y = y
		if last != nil {
			next := bar.Y - last.Y
			var char string
			currY := y
			switch {
			case next > 0:
				char = "╱"
				for i := x; i < x*spacing; i++ {
					currY--
					matrix[currY][x] = &api.Bar{Char: char}
				}
			case next < 0:
				char = "╲"
				for i := x; i < x*spacing; i++ {
					currY++
					matrix[currY][x] = &api.Bar{Char: char}
				}
			case next == 0:
				char = "─"
				for i := x; i < x*spacing; i++ {
					matrix[currY][x] = &api.Bar{Char: char}
				}
			}
		}
		last = bar
	}
	for _, slc := range matrix {
		out += "┃"
		for _, ptr := range slc {
			if ptr != nil {
				out += ptr.Char
			} else {
				out += " "
			}
		}
		out += "\n"
	}
	return out, nil
}
