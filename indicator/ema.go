package indicator

import (
	"anidiot.com/auto-trade/entity"
)

var K_9 = 0.2
var K_12 = 0.153846154
var K_25 = 0.076923077
var K_50 = 0.039215686
var K_200 = 0.009950249

func CurrentEMA(currentClose, previousEMA, k float64) float64 {
	return (k * (currentClose - previousEMA)) + previousEMA
}

func CurrentEMA12(currentClose float64, emaList entity.PrevEMA) float64 {
	return CurrentEMA(currentClose, emaList.EMA12, K_12)
}

func CurrentEMA25(currentClose float64, emaList entity.PrevEMA) float64 {
	return CurrentEMA(currentClose, emaList.EMA25, K_25)
}

func CurrentEMA50(currentClose float64, emaList entity.PrevEMA) float64 {
	return CurrentEMA(currentClose, emaList.EMA50, K_50)
}

func CalculateSMA(candles []entity.BinCandle) float64 {
	num := float64(len(candles))
	sum := 0.0
	for _, candle := range candles {
		sum += candle.Close
	}
	return sum / num
}

type emaIndicator struct {}

func (e emaIndicator) Calculate(oldData []entity.BinCandle, opt *entity.ChartOption) {
	ema := CalculateSMA(oldData[0:12:12])
	opt.PrevEMAList.EMA12=ema
	opt.PrevEMAList.EMA25=ema
	opt.PrevEMAList.EMA50=ema
	for i := 13; i < len(oldData); i++{
		oldData[i].EMA12 = CurrentEMA12(oldData[i].Close, opt.PrevEMAList)
		oldData[i].EMA25 = CurrentEMA25(oldData[i].Close, opt.PrevEMAList)
		oldData[i].EMA50 = CurrentEMA50(oldData[i].Close, opt.PrevEMAList)
		opt.PrevEMAList.EMA12=oldData[i].EMA12
		opt.PrevEMAList.EMA25=oldData[i].EMA25
		opt.PrevEMAList.EMA50=oldData[i].EMA50
	}
}

func (e emaIndicator) Prepare(candle *entity.BinCandle, oldData []entity.BinCandle, opt *entity.ChartOption) {
	prev := oldData[len(oldData)-1]
	opt.PrevEMAList.EMA12 =  prev.EMA12
	opt.PrevEMAList.EMA25 =  prev.EMA25
	opt.PrevEMAList.EMA50 =  prev.EMA50
}

func (e emaIndicator) CalculateCurrent(candle *entity.BinCandle, opt *entity.ChartOption) {
	candle.EMA12 = CurrentEMA12(candle.Close, opt.PrevEMAList)
	candle.EMA25 = CurrentEMA25(candle.Close, opt.PrevEMAList)
	candle.EMA50 = CurrentEMA50(candle.Close, opt.PrevEMAList)
}
