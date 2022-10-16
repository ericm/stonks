package graph

// ChartThemeType to change characters
type ChartThemeType int

const (
	// LineTheme is the lines chart theme
	LineTheme ChartThemeType = iota
	// DotTheme is the dots chart theme
	DotTheme
	// IconTheme is the icon chart theme
	IconTheme
)

// ChartTheme to change characters
type ChartTheme struct {
	Name     string
	UpChar   string
	FlatChar string
	DownChar string
}

// NewChartTheme returns the corresponding ChartTheme as indicated by themeType
func NewChartTheme(themeType ChartThemeType) ChartTheme {
	switch themeType {
	case DotTheme:
		return dotTheme
	case IconTheme:
		return iconTheme
	default:
		return lineTheme
	}
}

// lineTheme is the lines chart theme
var lineTheme ChartTheme = ChartTheme{
	Name:     "line",
	UpChar:   "╱",
	FlatChar: "─",
	DownChar: "╲",
}

// dotTheme is the dots chart theme
var dotTheme ChartTheme = ChartTheme{
	Name:     "dot",
	UpChar:   "·",
	FlatChar: "·",
	DownChar: "·",
}

// iconTheme is the icon chart theme
var iconTheme ChartTheme = ChartTheme{
	Name:     "icon",
	UpChar:   "⬆",
	FlatChar: "❚",
	DownChar: "⬇",
}
