/*
Package candlehandler provides handling of Japanese candles.
*/
package candlehandler

import (
	"encoding/csv"
	"fmt"
	"lection03/domain"
	"math"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type CandleMap = map[string]domain.Candle

// CandleHandler is a struct for handling Japanese candles.
type CandleHandler struct {
	Logger *logrus.Logger
}

// Handle executes the pipeline process.
func (h *CandleHandler) Handle(prices <-chan domain.Price) {
	wg := &sync.WaitGroup{}
	candles := pricesToCandles(prices, wg)

	for _, period := range [...]domain.CandlePeriod{
		domain.CandlePeriod1m,
		domain.CandlePeriod2m,
		domain.CandlePeriod10m,
	} {
		candles = h.convertCandles(candles, period, wg)
		candles = h.saveCandles(candles, period, wg)
	}

	for range candles {
	}

	wg.Wait()
}

// pricesToCandles converts domain.Price channel to domain.Candle channel, representing the price as a special case of
// a candle.
func pricesToCandles(prices <-chan domain.Price, wg *sync.WaitGroup) <-chan domain.Candle {
	candles := make(chan domain.Candle)

	wg.Add(1)
	go func() {
		defer func() {
			close(candles)
			wg.Done()
		}()

		for price := range prices {
			candles <- domain.Candle{
				Ticker: price.Ticker,
				Period: "",
				Open:   price.Value,
				High:   price.Value,
				Low:    price.Value,
				Close:  price.Value,
				TS:     price.TS,
			}
		}
	}()

	return candles
}

// convertCandles converts candles from current to given period.
func (h *CandleHandler) convertCandles(
	in <-chan domain.Candle, period domain.CandlePeriod, wg *sync.WaitGroup,
) <-chan domain.Candle {
	out := make(chan domain.Candle)

	wg.Add(1)
	go func() {
		defer func() {
			close(out)
			wg.Done()
		}()

		ts := time.Time{}
		candleMap := CandleMap{}

		closeCandles := func(candleMap CandleMap) {
			for _, candle := range candleMap {
				out <- candle
			}
		}

		for candleIn := range in {
			currentTS, err := domain.PeriodTS(period, candleIn.TS)

			if err != nil {
				h.Logger.Fatalln(err)
			}

			if !currentTS.Equal(ts) && len(candleMap) != 0 {
				closeCandles(candleMap)
				candleMap = CandleMap{}
			}

			ts = currentTS
			candle, ok := candleMap[candleIn.Ticker]

			if ok {
				candle = domain.Candle{
					Ticker: candle.Ticker,
					Period: candle.Period,
					Open:   candle.Open,
					High:   math.Max(candle.High, candleIn.High),
					Low:    math.Min(candle.Low, candleIn.Low),
					Close:  candleIn.Close,
					TS:     candle.TS,
				}
			} else {
				candle = domain.Candle{
					Ticker: candleIn.Ticker,
					Period: period,
					Open:   candleIn.Open,
					High:   candleIn.High,
					Low:    candleIn.Low,
					Close:  candleIn.Close,
					TS:     ts,
				}
			}

			candleMap[candleIn.Ticker] = candle
		}

		closeCandles(candleMap)
	}()

	return out
}

// saveCandles saves candles to the appropriate file.
func (h *CandleHandler) saveCandles(
	in <-chan domain.Candle, period domain.CandlePeriod, wg *sync.WaitGroup,
) <-chan domain.Candle {
	out := make(chan domain.Candle)
	file, err := os.Create(fmt.Sprintf("candles_%s.csv", period))

	if err != nil {
		h.Logger.Fatalln("file creation error")
	}

	writer := csv.NewWriter(file)

	wg.Add(1)
	go func() {
		defer func() {
			close(out)
			wg.Done()
		}()

		for candle := range in {
			fields := []string{
				candle.Ticker,
				fmt.Sprint(candle.TS),
				fmt.Sprint(candle.Open),
				fmt.Sprint(candle.High),
				fmt.Sprint(candle.Low),
				fmt.Sprint(candle.Close),
			}

			if err := writer.Write(fields); err != nil {
				h.Logger.Fatalln("error writing to csv:", err)
			}

			writer.Flush()

			if err := writer.Error(); err != nil {
				h.Logger.Fatalln(err)
			}

			out <- candle
		}
	}()

	return out
}
