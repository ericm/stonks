package graph

import "github.com/ericm/stonks/api"

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
	return out, nil
}
