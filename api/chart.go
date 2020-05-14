package api

import (
	"fmt"

	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
	"github.com/shopspring/decimal"
)

// Chart used to generate graphs
type Chart struct {
	Start    *datetime.Datetime
	End      *datetime.Datetime
	Length   int
	Interval datetime.Interval
	Bars     []*Bar
}

// Bar of a Chart
type Bar struct {
	Timestamp *datetime.Datetime
	Curremt   decimal.Decimal
}

// GetChart returns a Chart
func GetChart(symbol string, interval datetime.Interval, date *datetime.Datetime) (*Chart, error) {
	var q *chart.Iter
	if date != nil {
		q = chart.Get(&chart.Params{Symbol: "AMD", Interval: interval, Start: date})
	} else {
		q = chart.Get(&chart.Params{Symbol: "AMD", Interval: interval})
	}
	chart := &Chart{
		Interval: interval,
		Start:    datetime.FromUnix(q.Meta().CurrentTradingPeriod.Regular.Start),
		End:      datetime.FromUnix(q.Meta().CurrentTradingPeriod.Regular.End),
	}
	for q.Next() {
		bar := &Bar{datetime.FromUnix(q.Bar().Timestamp), q.Bar().Close}
		chart.Bars = append(chart.Bars, bar)
	}
	if len(chart.Bars) == 0 {
		return nil, fmt.Errorf("No bars were found for this time period")
	}
	return chart, nil
}
