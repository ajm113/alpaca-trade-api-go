package main

import (
	"fmt"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
)

func main() {
	// Get AAPL and MSFT trades from the first second of the 2021-08-09 market open
	multiTrades, err := marketdata.GetMultiTrades([]string{"AAPL", "MSFT"}, marketdata.GetTradesParams{
		Start: time.Date(2021, 8, 9, 13, 30, 0, 0, time.UTC),
		End:   time.Date(2021, 8, 9, 13, 30, 1, 0, time.UTC),
	})
	if err != nil {
		panic(err)
	}
	for symbol, trades := range multiTrades {
		fmt.Println(symbol + " trades:")
		for _, trade := range trades {
			fmt.Printf("%+v\n", trade)
		}
	}
	fmt.Println()

	// Get first 30 TSLA quotes from 2021-08-09 market open
	quotes, err := marketdata.GetQuotes("TSLA", marketdata.GetQuotesParams{
		Start:      time.Date(2021, 8, 9, 13, 30, 0, 0, time.UTC),
		TotalLimit: 30,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("TSLA quotes:")
	for _, quote := range quotes {
		fmt.Printf("%+v\n", quote)
	}
	fmt.Println()

	// Get all the IBM and GE 5-minute bars from the first half hour of the 2021-08-09 market open
	for item := range marketdata.GetMultiBarsAsync([]string{"IBM", "GE"}, marketdata.GetBarsParams{
		TimeFrame:  marketdata.NewTimeFrame(5, marketdata.Min),
		Adjustment: marketdata.Split,
		Start:      time.Date(2021, 8, 9, 13, 30, 0, 0, time.UTC),
		End:        time.Date(2021, 8, 9, 14, 0, 0, 0, time.UTC),
	}) {
		if err := item.Error; err != nil {
			panic(err)
		}
		fmt.Printf("%s: %+v\n", item.Symbol, item.Bar)
	}
	fmt.Println()

	// Get Facebook bars
	bars, err := marketdata.GetBars("META", marketdata.GetBarsParams{
		TimeFrame: marketdata.OneDay,
		Start:     time.Date(2022, 6, 1, 0, 0, 0, 0, time.UTC),
		End:       time.Date(2022, 6, 22, 0, 0, 0, 0, time.UTC),
		AsOf:      "2022-06-10", // Leaving it empty yields the same results
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("META bars:")
	for _, bar := range bars {
		fmt.Printf("%+v\n", bar)
	}
	fmt.Println()

	// Get Average Daily Trading Volume
	start := time.Date(2021, 8, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC)
	averageVolume, count, err := getADTV("AAPL", start, end)
	if err != nil {
		panic(err)
	}
	fmt.Printf("AAPL ADTV: %.2f (%d marketdays)\n", averageVolume, count)

	// Get news
	news, err := marketdata.GetNews(marketdata.GetNewsParams{
		Symbols:    []string{"AAPL", "TSLA"},
		Start:      time.Date(2021, 5, 6, 0, 0, 0, 0, time.UTC),
		End:        time.Date(2021, 5, 7, 0, 0, 0, 0, time.UTC),
		TotalLimit: 4,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("news:")
	for _, n := range news {
		fmt.Printf("%+v\n", n)
	}
}

func getADTV(symbol string, start, end time.Time) (av float64, n int, err error) {
	var (
		totalVolume uint64
	)
	for item := range marketdata.GetBarsAsync(symbol, marketdata.GetBarsParams{
		Start: start,
		End:   end,
	}) {
		if err = item.Error; err != nil {
			return
		}
		totalVolume += item.Bar.Volume
		n++
	}
	if n == 0 {
		return
	}
	av = float64(totalVolume) / float64(n)
	return
}
