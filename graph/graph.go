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
	spacing := (width - 5) / chart.Length
	out += "\n"
	var last *api.Bar
	for x, bar := range chart.Bars {
		y := int(bar.Current.Sub(chart.Low).Div(ran).Mul(
			decimal.NewFromInt((int64(height)))).Floor().IntPart())
		matrix[y][x*spacing] = bar
		bar.Y = y
		if last != nil {
			last.Next = bar.Y - last.Y
		}
		last = bar
	}
	for _, slc := range matrix {
		out += "┃"
		var last *api.Bar
		for _, ptr := range slc {
			if ptr != nil {
				out += "─"
				last = ptr
			} else {
				if last != nil {
					switch {
					case last.Next > 0:
						out += "╱"
					case last.Next < 0:
						out += "╲"
					case last.Next == 0:
						out += "─"
					}
				} else {
					out += " "
				}
			}
		}
		out += "\n"
	}
	return out, nil
}
