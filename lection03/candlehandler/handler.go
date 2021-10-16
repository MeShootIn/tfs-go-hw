/*
Package candlehandler provides work with Japanese candles.
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

type (
	CandleMap map[string]domain.Candle
	WriterMap map[domain.CandlePeriod]*csv.Writer
)

// CandleHandler is a struct for handling Japanese candles.
type CandleHandler struct {
	logger    *logrus.Logger
	writerMap WriterMap
}

// Config is an intermediary parameter for the NewCandleHandler constructor.
type Config struct {
	Logger *logrus.Logger
}

// NewCandleHandler is a public constructor for the CandleHandler struct.
func NewCandleHandler(config Config) (*CandleHandler, error) {
	getWriter := func(fileName string) (*csv.Writer, error) {
		file, err := os.Create(fileName)

		if err != nil {
			return nil, err
		}

		return csv.NewWriter(file), nil
	}

	candleHandler := CandleHandler{
		logger:    config.Logger,
		writerMap: WriterMap{},
	}

	for _, period := range [...]domain.CandlePeriod{
		domain.CandlePeriod1m,
		domain.CandlePeriod2m,
		domain.CandlePeriod10m,
	} {
		writer, err := getWriter("candles_" + string(period) + ".csv")

		if err != nil {
			return nil, err
		}

		candleHandler.writerMap[period] = writer
	}

	return &candleHandler, nil
}

// pricesToCandles converts domain.Price to domain.Candle channel, representing the price as a special case of a candle.
func pricesToCandles(wg *sync.WaitGroup, prices <-chan domain.Price) <-chan domain.Candle {
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

// Handle executes the pipeline process.
func (h *CandleHandler) Handle(prices <-chan domain.Price) {
	wg := &sync.WaitGroup{}
	candles := pricesToCandles(wg, prices)

	for _, period := range [...]domain.CandlePeriod{
		domain.CandlePeriod1m,
		domain.CandlePeriod2m,
		domain.CandlePeriod10m,
	} {
		candles = h.convertCandles(candles, wg, period)
	}

	for range candles {
	}

	wg.Wait()
}

// printCandleMap prints candles that are in the same period in .csv file.
func (h *CandleHandler) printCandleMap(candleMap CandleMap, period domain.CandlePeriod) {
	for _, candle := range candleMap {
		fields := []string{
			candle.Ticker,
			fmt.Sprint(candle.TS),
			fmt.Sprint(candle.Open),
			fmt.Sprint(candle.High),
			fmt.Sprint(candle.Low),
			fmt.Sprint(candle.Close),
		}

		if err := h.writerMap[period].Write(fields); err != nil {
			h.logger.Fatalln("error writing to csv:", err)
		}
	}

	h.writerMap[period].Flush()

	if err := h.writerMap[period].Error(); err != nil {
		h.logger.Fatalln(err)
	}
}

// convertCandles function converts candles from current to given period.
func (h *CandleHandler) convertCandles(
	in <-chan domain.Candle, wg *sync.WaitGroup, period domain.CandlePeriod,
) <-chan domain.Candle {
	out := make(chan domain.Candle)

	handleCandlesTS := func(candleMap CandleMap) {
		h.printCandleMap(candleMap, period)

		for _, candle := range candleMap {
			out <- candle
		}
	}

	wg.Add(1)
	go func() {
		defer func() {
			close(out)
			wg.Done()
		}()

		ts := time.Time{}
		candleMap := CandleMap{}

		for candleIn := range in {
			periodTS, err := domain.PeriodTS(period, candleIn.TS)

			if err != nil {
				h.logger.Errorln(err)
			}

			if !periodTS.Equal(ts) && len(candleMap) != 0 {
				handleCandlesTS(candleMap)
				candleMap = CandleMap{}
			}

			ts = periodTS
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

		handleCandlesTS(candleMap)
	}()

	return out
}
