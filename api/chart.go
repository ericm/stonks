package api

import (
	"fmt"

	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
	"github.com/shopspring/decimal"
)

// Chart used to generate graphs
type Chart struct {
	Ticker   string
	Exchange string
	Currency string
	Start    *datetime.Datetime
	End      *datetime.Datetime
	Length   int
	High     decimal.Decimal
	Low      decimal.Decimal
	Open     decimal.Decimal
	Close    decimal.Decimal
	Interval datetime.Interval
	Bars     []*Bar
	Change   decimal.Decimal
}

// Bar of a Chart
type Bar struct {
	Timestamp *datetime.Datetime
	Current   decimal.Decimal
	Y         int
	Char      string
}

// GetChart returns a Chart
func GetChart(symbol string, interval datetime.Interval, start *datetime.Datetime, end *datetime.Datetime) (*Chart, error) {
	q := chart.Get(&chart.Params{Symbol: symbol, Interval: interval, Start: start, End: end, IncludeExt: false})
	if q.Count() < 7 && interval == datetime.FifteenMins {
		q = chart.Get(&chart.Params{Symbol: symbol, Interval: datetime.FiveMins, Start: start, End: end, IncludeExt: false})
	}
	var chart *Chart
	for q.Next() {
		if chart == nil {
			chart = &Chart{
				Interval: interval,
				Start:    datetime.FromUnix(q.Bar().Timestamp),
				High:     q.Bar().High,
				Low:      q.Bar().Low,
				Open:     q.Bar().Open,
				Length:   q.Count(),
				Ticker:   symbol,
				Exchange: q.Meta().ExchangeName,
				Currency: q.Meta().Currency,
			}
		}
		bar := &Bar{Timestamp: datetime.FromUnix(q.Bar().Timestamp), Current: q.Bar().Close}
		if bar.Current.IsZero() {
			continue
		}
		if q.Bar().High.GreaterThan(chart.High) {
			chart.High = q.Bar().High
		}
		if q.Bar().Low.LessThan(chart.Low) {
			chart.Low = q.Bar().Low
		}
		chart.Close = q.Bar().Close
		chart.End = datetime.FromUnix(q.Bar().Timestamp)
		chart.Bars = append(chart.Bars, bar)
	}
	if chart == nil || len(chart.Bars) == 0 {
		return nil, fmt.Errorf("No bars were found for this time period for %s", symbol)
	}
	chart.Change = chart.Close.Sub(chart.Open).Div(chart.Open).Mul(decimal.NewFromInt(100))
	return chart, nil
}
