package graph

import (
	"fmt"
	"strings"
	"time"

	"github.com/ericm/stonks/api"
	"github.com/shopspring/decimal"
)

const (
	dateFormat = "Mon 02/01/2006 15:04"
	timeFormat = "3.04pm"
	dayFormat  = "2 Jan"
)

// ChartTheme to change characters
type ChartTheme int

const (
	// LineTheme is the lines chart theme
	LineTheme ChartTheme = iota
	// DotTheme is the dots chart theme
	DotTheme
	// IconTheme is the icon chart theme
	IconTheme
)

func borderHorizontal(width int) string {
	return strings.Repeat("━", width)
}

// infoHeader builds the contents of the top part of the graph,
// where ticker details are found
func infoHeader(chart *api.Chart, width int, maxSize int) string {
	out := "┏"
	out += borderHorizontal(width + maxSize + 1)
	out += "┓"

	colour := 92
	if chart.Change.IsNegative() {
		colour = 91
	}

	info := fmt.Sprintf(
		"\n┃\033[95m %5s | \033[%dm%s %s (%s%% | %s)\033[95m on %s | ",
		chart.Ticker,
		colour,
		chart.Close.StringFixed(2),
		chart.Currency,
		chart.Change.StringFixed(2),
		chart.ChangeVal.StringFixed(2),
		chart.End.Time().Format(dateFormat),
	)
	if len(info) > width {
		info += fmt.Sprintf(
			"%s \033[0m",
			chart.Exchange,
		)
	} else {
		info += fmt.Sprintf(
			"Prev: %s | %s \033[0m",
			chart.Prev.StringFixed(2),
			chart.Exchange,
		)
	}
check:
	if len(info) < width+maxSize+24 {
		info += " "
		goto check
	}
	info += "┃\n┣"
	out += info
	out += borderHorizontal(width + maxSize + 1)
	out += "┫\n"

	return out
}

// chartArea renders the different bars in the chart along with
// relevant price levels and returns it as a string
func chartArea(chart *api.Chart, width int, height int, maxSize int, chartTheme ChartTheme) string {
	if chart.Length < width/5 {
		chart.Length = width / 3
	}
	if chart.Length > width {
		mod := 2 * (chart.Length / width)
		chartTemp := make([]*api.Bar, 0)
		for i, bar := range chart.Bars {
			if (i+1)%mod == 0 {
				chartTemp = append(chartTemp, bar)
			}
		}
		chart.Bars = chartTemp
		chart.Length = len(chartTemp)
	}

	matrix := make([][]*api.Bar, height)
	for i := range matrix {
		matrix[i] = make([]*api.Bar, width)
	}
	ran := chart.High.Sub(chart.Low)
	spacing := (width) / (chart.Length)
	if spacing == 0 {
		spacing = 3
	}

	var last *api.Bar
	var (
		upChar   = "╱"
		flatChar = "─"
		downChar = "╲"
	)
	if chartTheme == DotTheme {
		upChar = "·"
		flatChar = "·"
		downChar = "·"
	} else if chartTheme == IconTheme {
		upChar = "⬆"
		flatChar = "❚"
		downChar = "⬇"
	}
	for x, bar := range chart.Bars {
		bar.Char = flatChar
		y := height - int(bar.Current.Sub(chart.Low).Div(ran).Mul(
			decimal.NewFromInt((int64(height)))).Floor().IntPart())
		if y >= height {
			y--
		}
		newX := x * spacing
		if newX >= width {
			newX = width - 1
			last = nil
		}
		matrix[y][newX] = bar
		bar.Y = y
		if last != nil {
			next := last.Y - bar.Y
			var char string
			currY := last.Y
			switch {
			case next > 0:
				char = upChar
				bar.Char = char
				for i := 0; i < spacing-1; i++ {
					currY--
					if currY >= 0 && currY >= y {
						matrix[currY][i+((x-1)*spacing)+1] = &api.Bar{Char: char}
					}
				}
			case next < 0:
				char = downChar
				bar.Char = char
				for i := 0; i < spacing-1; i++ {
					currY++
					if currY < height && currY <= y {
						matrix[currY][i+((x-1)*spacing)+1] = &api.Bar{Char: char}
					}
				}
			case next == 0:
				char = flatChar
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

	out := ""
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

	return out
}

// timeAxisFooter builds the bottom part of the graph, where time marks
// indicate when the prices took place, and returns it as a string
func timeAxisFooter(chart *api.Chart, width int, maxSize int, timezone *time.Location) string {
	out := "┣"
	out += borderHorizontal(width + maxSize + 1)
	out += "┫\n"

	footer := "┃"
incFooter:
	if len(footer) < maxSize+4 {
		footer += " "
		goto incFooter
	}
	mod := width / chart.Length
	if mod < 3 {
		mod = width / 10
	}

	spacing := (width) / (chart.Length)
	if spacing == 0 {
		spacing = 3
	}

retryFooter:
	diff := mod * spacing
	lastLen := 0
	for i, bar := range chart.Bars {
		if i%mod == 0 {
			format := timeFormat
			if chart.End.Unix()-chart.Start.Unix() > 86400 {
				format = dayFormat
			}
			t := bar.Timestamp.Time().In(timezone).Format(format)
			if lastLen > 0 {
				for _i := 0; _i < diff-len(t); _i++ {
					footer += " "
				}
			}
			footer += t
			lastLen = len(t)
		}
	}
	if len(footer) > width+maxSize+4 {
		mod++
		goto retryFooter
	}
checkFooter:
	if len(footer) < width+maxSize+4 {
		footer += " "
		goto checkFooter
	}
	footer += "┃"
	out += footer
	out += "\n┗"
	out += borderHorizontal(width + maxSize + 1)
	out += "┛\n"

	return out
}

// GenerateGraph with ASCII graph with ANSI escapes
func GenerateGraph(chart *api.Chart, width int, height int, chartTheme ChartTheme, timezone *time.Location) (string, error) {
	maxSize := len(strings.Split(chart.High.String(), ".")[0]) + 3 // Add 3 for the dot and precision 2.
	if maxSize < 7 {
		maxSize += 7 % maxSize
	}

	out := infoHeader(chart, width, maxSize)
	out += chartArea(chart, width, height, maxSize, chartTheme)
	out += timeAxisFooter(chart, width, maxSize, timezone)

	return out, nil
}
