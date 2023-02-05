package indicator

import (
	"anidiot.com/auto-trade/entity"
	"anidiot.com/auto-trade/utils"
)

type HeikinAshiCandle struct {}

func (c HeikinAshiCandle) Prepare(candle *entity.BinCandle, oldData []entity.BinCandle, opt *entity.ChartOption) {
	previous := oldData[len(oldData)-1]
	candle.HAOpen = utils.AverageFloat(previous.HAOpen, previous.HAClose)
	candle.HAClose = utils.AverageFloat(candle.Open, candle.Close,
		candle.High, candle.Low)
	candle.HAHigh, candle.HALow = utils.FindMaxNMinFloat(candle.High, candle.Low, candle.HAOpen, candle.HAClose)
}

func (c HeikinAshiCandle) CalculateCurrent(candle *entity.BinCandle, opt *entity.ChartOption) {
	candle.HAClose = utils.AverageFloat(candle.Open, candle.Close,
		candle.High, candle.Low)
	candle.HAHigh, candle.HALow = utils.FindMaxNMinFloat(candle.High,
		candle.Low, candle.HAOpen, candle.HAClose)
}

func (HeikinAshiCandle) Calculate(candles []entity.BinCandle, opt *entity.ChartOption) {
	candles[0].HAClose = utils.AverageFloat(candles[0].Open, candles[0].Close, candles[0].High, candles[0].Low)
	candles[0].HAOpen = candles[0].Open
	candles[0].HALow =  candles[0].Low
	candles[0].HAHigh = candles[0].High
	for i:=1; i < len(candles); i++ {
		previous := candles[i-1]
		tick := candles[i]
		tick.HAOpen = utils.AverageFloat(previous.Open, previous.Close)
		tick.HAClose = utils.AverageFloat(tick.Open, tick.Close,
			tick.High, tick.Low)
		tick.HAHigh, tick.HALow =
			utils.FindMaxNMinFloat(tick.High, tick.Low, tick.HAOpen, tick.HAClose)
	}
}
