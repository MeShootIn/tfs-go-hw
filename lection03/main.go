package main

import (
	"context"
	"lection03/candlehandler"
	"lection03/generator"
	"os"
	"os/signal"
	"time"

	"github.com/sirupsen/logrus"
)

var tickers = []string{"AAPL", "SBER", "NVDA", "TSLA"}

func main() {
	logger := logrus.New()
	ctx, cancel := context.WithCancel(context.Background())

	// Signal channel for catching SIGINT (Ctrl+C) signals.
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	defer func() {
		signal.Stop(signalChan)
		close(signalChan)
	}()

	pg := generator.NewPricesGenerator(generator.Config{
		Factor:  10,
		Delay:   time.Millisecond * 500,
		Tickers: tickers,
	})
	logger.Info("start prices generator...")
	prices := pg.Prices(ctx)

	ch, err := candlehandler.NewCandleHandler(candlehandler.Config{
		Logger: logger,
	})

	if err != nil {
		logger.Fatalln(err)

		return
	}

	// Catching SIGINT and cancelling the ctx context
	go func() {
		<-signalChan
		cancel()
		logger.Info("price generation process cancelled")
	}()

	// Main handling process
	ch.Handle(prices)
	logger.Info("candles for all periods are saved in the candles_{period}.csv files")
}
