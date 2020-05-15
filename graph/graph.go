package graph

import (
	"time"

	"github.com/ericm/stonks/api"
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
	interval, err := time.ParseDuration(string(chart.Interval))
	if err != nil {
		return "", err
	}
	count := interval.Seconds() / float64(width)
	difference := int(interval.Seconds()) % width

	curr := chart.Bars[0]
	for y, bar := range chart.Bars[1:] {

	}
	return out, nil
}
