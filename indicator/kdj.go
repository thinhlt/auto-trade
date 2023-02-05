package indicator

import (
	"anidiot.com/auto-trade/entity"
	"math"
)

type kdjIndicator struct {}

func (indicator kdjIndicator) Calculate(oldData []entity.BinCandle, opt *entity.ChartOption) {
	oldData[8].KDJ = entity.KDJ{K: 50, D: 50}
	for i := 9; i < len(oldData); i++{
		kdj, max, min := prepare(oldData[i-9: i : 9])
		prev := oldData[i-1]

		rvs := (prev.Close - min) / (max - min) * 100
		oldData[i].K = ((2.0 / 3) * kdj.K) + ((1.0 / 3) * rvs)
		oldData[i].D = ((2.0 / 3) * kdj.D) + ((1.0 / 3) * oldData[i].K)
		oldData[i].J = (3 * oldData[i].K) - (2 * oldData[i].D)
	}
}

func (kdjIndicator) Prepare(candle *entity.BinCandle, oldData []entity.BinCandle, opt *entity.ChartOption) {
	kdj, max, min := prepare(oldData[len(oldData)-9:])
	opt.High9 = max
	opt.Low9 = min
	opt.PrevKDJ = kdj
}

func (kdjIndicator) CalculateCurrent(candle *entity.BinCandle, opt *entity.ChartOption) {
	rvs := (opt.PrevClose - opt.Low9) / (opt.High9 - opt.Low9) * 100
	candle.K = ((2.0 / 3) * opt.PrevKDJ.K) + ((1.0 / 3) * rvs)
	candle.D = ((2.0 / 3) * opt.PrevKDJ.D) + ((1.0 / 3) * opt.PrevKDJ.K)
	candle.J = (3 * candle.K) - (2 * candle.D)
}

func prepare(oldData []entity.BinCandle) (kdj entity.KDJ, max, min float64) {
	prev := oldData[len(oldData)-1]
	min = math.MaxFloat64
	max = float64(0)
	for i := 0; i < len(oldData); i++ {
		max = math.Max(max, oldData[i].High)
		min = math.Min(min, oldData[i].Low)
	}
	kdj = prev.KDJ
	return
}
