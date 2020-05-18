package graph

import (
	"fmt"
	"strings"

	"github.com/ericm/stonks/api"
	"github.com/shopspring/decimal"
)

const dateFormat = "Mon 02/01/2006 15:04 GMT"
const timeFormat = "3.04pm"

func borderHorizontal(out *string, width int) {
	for _i := 0; _i < width-2; _i++ {
		*out += "━"
	}
}

// GenerateGraph with ASCII graph with ANSI escapes
func GenerateGraph(chart *api.Chart, width int, height int) (string, error) {
	out := "┏"
	maxSize := len(strings.Split(chart.High.String(), ".")[0]) + 3
	borderHorizontal(&out, width+maxSize+3)
	out += "┓"
	info := fmt.Sprintf(
		"\n┃\033[95m %s - \033[92m%s %s\033[95m on %s - %s \033[0m",
		chart.Ticker,
		chart.Close.StringFixed(2),
		chart.Currency,
		chart.End.Time().Format(dateFormat),
		chart.Exchange,
	)
check:
	if len(info) < width+maxSize+24 {
		info += " "
		goto check
	}
	info += "┃\n┣"
	out += info
	borderHorizontal(&out, width+maxSize+3)
	out += "┫"
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
		y := height - int(bar.Current.Sub(chart.Low).Div(ran).Mul(
			decimal.NewFromInt((int64(height)))).Floor().IntPart())
		if y >= height {
			y--
		}
		newX := x * spacing
		if newX >= width {
			newX--
		}
		matrix[y][newX] = bar
		bar.Y = y
		if last != nil {
			next := last.Y - bar.Y
			var char string
			currY := last.Y
			switch {
			case next > 0:
				char = "╱"
				bar.Char = char
				for i := 0; i < spacing-1; i++ {
					currY--
					if currY >= 0 && currY >= y {
						matrix[currY][i+((x-1)*spacing)+1] = &api.Bar{Char: char}
					}
				}
			case next < 0:
				char = "╲"
				bar.Char = char
				for i := 0; i < spacing-1; i++ {
					currY++
					if currY < height && currY <= y {
						matrix[currY][i+((x-1)*spacing)+1] = &api.Bar{Char: char}
					}
				}
			case next == 0:
				char = "─"
				last.Char = char
				for i := 0; i < spacing-1; i++ {
					matrix[currY][i+((x-1)*spacing)+1] = &api.Bar{Char: char}
				}
			}
			// Edge cases
			switch last.Char {
			case "╱":
				switch char {
				case "╲":
					if newX > 0 && matrix[y][(newX)-1] != nil {
						last.Char = "▁"
					} else {
						last.Char = "ʌ"
					}
				case "╱":
					last.Char = "╱"
				}
			case "╲":
				switch char {
				case "╲":
					if newX > 0 && matrix[y][(newX)-1] != nil {
						last.Char = "▔"
					} else {
						last.Char = "▁"
					}
				case "╱":
					last.Char = "╱"
				}
			}
		}
		last = bar
	}
	increment := ran.Div(decimal.NewFromInt(int64(height)))
	for i, slc := range matrix {
		out += "┃"
		price := chart.High.Sub(increment.Mul(decimal.NewFromInt(int64(i)))).StringFixed(2)
	checkLen:
		if len(price) < maxSize {
			price = " " + price
			goto checkLen
		}
		out += price
		out += "│\033[92m"
		for _, ptr := range slc {
			if ptr != nil {
				out += ptr.Char
			} else {
				out += " "
			}
		}
		out += "\033[0m┃"
		out += "\n"
	}
	out += "┣"
	borderHorizontal(&out, width+maxSize+3)
	out += "┫\n"
	footer := "┃"
incFooter:
	if len(footer) < maxSize+4 {
		footer += " "
		goto incFooter
	}
	mod := width / chart.Length
	if mod == 1 {
		mod = width / 10
	}
	fmt.Println(mod, chart.Length)
	diff := mod * spacing
	lastLen := 0
	for i, bar := range chart.Bars {
		if i%mod == 0 {
			t := bar.Timestamp.Time().Format(timeFormat)
			if lastLen > 0 {
				for _i := 0; _i < diff-len(t); _i++ {
					footer += " "
				}
			}
			footer += t
			lastLen = len(t)
		}
	}
checkFooter:
	if len(footer) < width+maxSize+4 {
		footer += " "
		goto checkFooter
	}
	footer += "┃"
	out += footer
	out += "\n┗"
	borderHorizontal(&out, width+maxSize+3)
	out += "┛\n"
	return out, nil
}
