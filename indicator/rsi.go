package indicator

import (
	"anidiot.com/auto-trade/entity"
	"math"
)

type rsiIndicator struct {}

func (rsiIndicator) Calculate(oldData []entity.BinCandle, opt *entity.ChartOption) {
	for i := 0; i < 16; i++ {
		oldData[i].RSI = 0
	}
	for i := 15; i < len(oldData); i++ {
		subSlice := oldData[i-15 : i+1]
		oldData[i].RSI = CurrentBinRSI(subSlice)
	}
}
func (rsiIndicator) Prepare(candle *entity.BinCandle, oldData []entity.BinCandle, opt *entity.ChartOption) {
	candle.RSI = CurrentBinRSI(oldData)
}

func (rsiIndicator) CalculateCurrent(candle *entity.BinCandle, opt *entity.ChartOption) {
	//candle.RSI = CurrentBinRSI(oldData)
}

func CurrentBinRSI(candles []entity.BinCandle) float64 {
	upChange, downChange := 0.0, 0.0
	for i := 1; i < len(candles); i++ {
		diff := math.Abs(candles[i].Close - candles[i-1].Close)
		if candles[i].Close > candles[i-1].Close {
			upChange += diff
		} else {
			downChange += diff
		}
	}
	rsi := 100 - 100/(1+(upChange/downChange))
	if upChange == 0 {
		return 0
	}
	if downChange == 0 {
		return 100
	}
	return rsi
}
