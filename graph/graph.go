package graph

import (
	"fmt"
	"strings"
	"time"

	"github.com/ericm/stonks/api"
	"github.com/exrook/drawille-go"
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
	// Braille theme
	Braille
)

func borderHorizontal(out *string, width int) {
	for _i := 0; _i < width-2; _i++ {
		*out += "━"
	}
}

func GenerateBraille(chart *api.Chart, width int, height int) string {

	ran := chart.High.Sub(chart.Low)
	wif := 159
	yScale := 4.2

	xScale := float64(wif) / float64(chart.Length)

	println(chart.Length)

	canvas := drawille.NewCanvas()

	lastX := -1.0
	lastY := -1.0

	for x, bar := range chart.Bars {

		xf := float64(x)

		y := height - int(bar.Current.Sub(chart.Low).Div(ran).Mul(
			decimal.NewFromInt((int64(height)))).Floor().IntPart())

		yf := float64(y)

		if lastX != -1 {

			canvas.DrawLine(lastX*xScale, lastY*yScale, xf*xScale, yf*yScale)

		}

		lastX = xf
		lastY = yf

	}

	return canvas.String()
}

func generateMatrix(chart *api.Chart, width int, height int, chartTheme ChartTheme, spacing int) [][]*api.Bar {
	ran := chart.High.Sub(chart.Low)
	matrix := make([][]*api.Bar, height)
	for i := range matrix {
		matrix[i] = make([]*api.Bar, width)
	}

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

	return matrix
}

// GenerateGraph with ASCII graph with ANSI escapes
func GenerateGraph(chart *api.Chart, width int, height int, chartTheme ChartTheme, timezone *time.Location) (string, error) {
	ran := chart.High.Sub(chart.Low)

	maxSize := len(strings.Split(chart.High.String(), ".")[0]) + 3

	spacing := (width) / (chart.Length)

	if spacing == 0 {
		spacing = 3
	}

	out := "┏"

	borderHorizontal(&out, width+maxSize+3)
	out += "┓"
	colour := 92
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
	if chart.Change.IsNegative() {
		colour = 91
	}
	info := fmt.Sprintf(
		"\n┃\033[95m %s | \033[%dm%s %s (%s%% | %s)\033[95m on %s | ",
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
	borderHorizontal(&out, width+maxSize+3)

	out += "┫"

	increment := ran.Div(decimal.NewFromInt(int64(height)))
	out += "\n"

	if chartTheme == Braille {
		splitCanvas := strings.Split(GenerateBraille(chart, width, height), "\n")

		for i, split := range splitCanvas {

			if len(split) == 0 {
				break
			}

			out += "┃"
			price := chart.High.Sub(increment.Mul(decimal.NewFromInt(int64(i)))).StringFixed(2)
		checkLen:
			if len(price) < maxSize {
				price = " " + price
				goto checkLen
			}
			out += price
			out += "│\033[92m"

			for _, ch := range split {

				out += string(ch)
			}

			out += "\033[0m┃"
			out += "\n"
		}
	} else {

		matrix := generateMatrix(chart, width, height, chartTheme, spacing)

		for i, slc := range matrix {
			out += "┃"
			price := chart.High.Sub(increment.Mul(decimal.NewFromInt(int64(i)))).StringFixed(2)
		checkLen2:
			if len(price) < maxSize {
				price = " " + price
				goto checkLen2
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

	}

	//

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
	if mod < 3 {
		mod = width / 10
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
	borderHorizontal(&out, width+maxSize+3)
	out += "┛\n"
	return out, nil

}
