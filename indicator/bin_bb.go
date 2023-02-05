package indicator

import (
	"math"

	"anidiot.com/auto-trade/entity"
)

type bbIndicator struct {}
func (bbIndicator) Calculate(oldData []entity.BinCandle, opt *entity.ChartOption) {
	closePrice := oldData[0].Close
	oldData[0].BBB = entity.BoilingerBand{
		Basic: closePrice,
		Upper: closePrice,
		Lower: closePrice,
	}
	closePrice = oldData[1].Close
	oldData[1].BBB = entity.BoilingerBand{
		Basic: closePrice,
		Upper: closePrice,
		Lower: closePrice,
	}
	i := 2
	for ; i < 21; i++ {
		preCandle := oldData[0:i]
		oldData[i].BBB = CurrentBBand(preCandle)
	}
	for ; i < len(oldData); i++ {
		start := i - 20
		preCandle := oldData[start:i]
		oldData[i].BBB = CurrentBBand(preCandle)
	}
}
func (bbIndicator) Prepare(candle *entity.BinCandle, oldData []entity.BinCandle, opt *entity.ChartOption) {

}
func (bbIndicator) CalculateCurrent(candle *entity.BinCandle, opt *entity.ChartOption) {

}

func CurrentBBand(candles []entity.BinCandle) entity.BoilingerBand {
	sum := float64(0)
	length := float64(len(candles))
	for i := range candles {
		sum += candles[i].Close
	}
	mean := sum / length
	v := 0.0
	for i := range candles {
		dif := candles[i].Close - mean
		v += math.Pow(dif, 2)
	}
	variance := math.Sqrt(v / (length - 1))
	return entity.BoilingerBand{
		Basic: mean,
		Upper: mean + 2*variance,
		Lower: mean - 2*variance,
	}
}

func CalculateBB(data []float64, ticker *entity.BinCandle, opts *entity.ChartOption) {
	sum := float64(0)
	length := float64(len(data))
	for i := range data {
		sum += data[i]
	}
	mean := sum / length
	v := 0.0
	for i := range data {
		dif := data[i] - mean
		v += math.Pow(dif, 2)
	}
	deviation := math.Sqrt(v / (length - 1))
	opts.Deviation = deviation
	opts.DeviationPercent = deviation / mean * 100
	ticker.BBB.Basic = mean
	ticker.BBB.Lower = mean - 2*deviation
	ticker.BBB.Upper = mean + 2*deviation
	opts.BBWidth = (100 * 4 * deviation) / ticker.BBB.Lower
}

