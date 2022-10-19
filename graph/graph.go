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

// GenerateGraph with ASCII graph with ANSI escapes
func GenerateGraph(chart *api.Chart, width int, height int, chartTheme ChartTheme, timezone *time.Location) (string, error) {
	maxSize := len(strings.Split(chart.High.String(), ".")[0]) + 3 // Add 3 for the dot and precision 2.
	if maxSize < 7 {
		maxSize += 7 % maxSize
	}

	var out strings.Builder
	out.WriteString(infoHeader(chart, width, maxSize))
	out.WriteString(chartArea(chart, width, height, maxSize, chartTheme))
	out.WriteString(timeAxisFooter(chart, width, maxSize, timezone))

	return out.String(), nil
}

// infoHeader builds the contents of the top part of the graph,
// where ticker details are found
func infoHeader(chart *api.Chart, width int, maxSize int) string {
	colour := 92
	if chart.Change.IsNegative() {
		colour = 91
	}

	info := fmt.Sprintf(
		"\033[95m %5s |\033[%dm %s %s (%s%% | %s)\033[95m on %s | ",
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

	// add padding to info section so that it fills the total graph width
	// total width = first column width + column separator (|) + width + 19 special formatting chars (not displayed)
	info = padString(info, maxSize+1+width+19, false)

	var out strings.Builder
	out.WriteString("┏" + borderHorizontal(width+maxSize+1) + "┓\n")
	out.WriteString("┃" + info + "┃\n")
	out.WriteString("┣" + borderHorizontal(width+maxSize+1) + "┫\n")
	return out.String()
}

// chartArea renders the different bars in the chart along with
// relevant price levels and returns it as a string
func chartArea(chart *api.Chart, width int, height int, maxSize int, chartTheme ChartTheme) string {
	if chart.Length < width/5 {
		chart.Length = width / 3
	}

	// sample bars in the chart if they all don't fit in the graph's width
	if chart.Length > width {
		mod := 2 * (chart.Length / width)
		barsTemp := make([]*api.Bar, 0)
		for i, bar := range chart.Bars {
			if i%mod == 0 {
				barsTemp = append(barsTemp, bar)
			}
		}
		chart.Bars = barsTemp
		chart.Length = len(barsTemp)
	}

	matrix := make([][]*api.Bar, height)
	for i := range matrix {
		matrix[i] = make([]*api.Bar, width)
	}

	ran := chart.High.Sub(chart.Low)
	spacing := width / chart.Length
	if spacing == 0 {
		spacing = 3
	}

	var last *api.Bar
	for x, bar := range chart.Bars {
		bar.Char = chartTheme.FlatChar
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
			var char string
			currY := last.Y
			switch {
			case last.Y > bar.Y:
				char = chartTheme.UpChar
				bar.Char = char
				for i := 0; i < spacing-1; i++ {
					currY--
					if currY >= 0 && currY >= y {
						matrix[currY][i+((x-1)*spacing)+1] = &api.Bar{Char: char}
					}
				}
			case last.Y < bar.Y:
				char = chartTheme.DownChar
				bar.Char = char
				for i := 0; i < spacing-1; i++ {
					currY++
					if currY < height && currY <= y {
						matrix[currY][i+((x-1)*spacing)+1] = &api.Bar{Char: char}
					}
				}
			case last.Y == bar.Y:
				char = chartTheme.FlatChar
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

	var out strings.Builder
	for i, line := range matrix {
		price := chart.High.Sub(increment.Mul(decimal.NewFromInt(int64(i)))).StringFixed(2)

		// left-pad the price column to right-justify prices
		price = padString(price, maxSize, true)

		out.WriteString("┃")
		out.WriteString(price)
		out.WriteString("│\033[92m")

		for _, char := range line {
			if char != nil {
				out.WriteString(char.Char)
			} else {
				out.WriteString(" ")
			}
		}

		out.WriteString("\033[0m┃\n")
	}

	return out.String()
}

// timeAxisFooter builds the bottom part of the graph, where time marks
// indicate when the prices took place, and returns it as a string
func timeAxisFooter(chart *api.Chart, width int, maxSize int, timezone *time.Location) string {
	mod := width / chart.Length
	if mod < 3 {
		mod = width / 10
	}

	spacing := width / chart.Length
	if spacing == 0 {
		spacing = 3
	}

	diff := mod * spacing

	format := timeFormat
	if chart.End.Unix()-chart.Start.Unix() > 86400 {
		format = dayFormat
	}

	// retry building the footer until a valid one (i.e. one that fits at the bottom of
	// the graph) is obtained. If the footer is not valid, mod will be increased by 1
	// and building it will be retried
	footer := ""
	for {
		footer = strings.Repeat(" ", maxSize+1)
		for i, bar := range chart.Bars {
			if i%mod == 0 {
				t := bar.Timestamp.Time().In(timezone).Format(format)
				t = padString(t, diff, false)

				footer += t
			}
		}

		footer = strings.TrimRight(footer, " ")

		if len(footer) <= maxSize+1+width {
			break
		}

		mod++
	}

	footer = padString(footer, maxSize+1+width, false)

	var out strings.Builder
	out.WriteString("┣" + borderHorizontal(width+maxSize+1) + "┫\n")
	out.WriteString("┃" + footer + "┃\n")
	out.WriteString("┗" + borderHorizontal(width+maxSize+1) + "┛\n")

	return out.String()
}

func borderHorizontal(width int) string {
	return strings.Repeat("━", width)
}

// padStrings pads a string with spaces until it fills the given width.
// If left is true, spaces are added at the beginning of the string
func padString(s string, width int, left bool) string {
	padFmt := fmt.Sprintf("%%-%ds", width)
	if left {
		padFmt = fmt.Sprintf("%%%ds", width)
	}

	return fmt.Sprintf(padFmt, s)
}
