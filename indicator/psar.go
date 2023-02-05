package indicator

import (
	"anidiot.com/auto-trade/entity"
)

const MAX_AF = 0.2
type psarIndicator struct {}

func (p psarIndicator) Calculate(candles []entity.BinCandle, opt *entity.ChartOption) {
	var hight, low, ep float64
	af, maxAF := 0.02, 0.2
	for i := range candles {
		switch i {
		case 0:
			low = candles[0].Low
			break
		case 1:
			if candles[0].Psar == 0 {
				if candles[1].Close > candles[0].Close {
					candles[1].Psar = candles[0].Low
					if candles[i].High > hight {
						ep = candles[i].High
					} else {
						ep = hight
					}
				} else {
					candles[1].Psar = candles[0].High
					if candles[i].Low < low {
						ep = candles[i].Low
					} else {
						ep = low
					}
				}
				break
			}
			hight = candles[i].High
			if candles[0].Psar < candles[0].Close {
				ep = hight
			} else {
				ep = low
			}
			fallthrough
		default:
			candle := candles[i]
			lastCandle := candles[i-1]
			prevSAR := lastCandle.Psar
			// down trend
			if lastCandle.Close < lastCandle.Psar {
				sar := prevSAR - af*(prevSAR-ep)
				// if current high > sar, set prev low is sar
				if candle.High > sar {
					candles[i].Psar = low
					hight = candle.High
					ep = hight
					af = 0.02
					break
				}
				// if new high, update EP, AF and high
				candles[i].Psar = sar
				if candle.Low < low {
					ep = candle.Low
					if af < maxAF {
						af += 0.02
					}
					break
				}
				// uptrend
			} else {
				sar := prevSAR + af*(ep-prevSAR)
				// if current low < sar, set prev high is sar
				if candle.Low < sar {
					candles[i].Psar = hight
					low = candle.Low
					ep = low
					af = 0.02
					break
				}
				// if new low, update EP, AF and low
				candles[i].Psar = sar
				if candle.High > hight {
					ep = candle.High
					if af < maxAF {
						af += 0.02
					}
				}
			}
		}

		if candles[i].Low < low {
			low = candles[i].Low
		}
		if candles[i].High > hight {
			hight = candles[i].High
		}
	}

	opt.High=hight
	opt.Low=low
	opt.AF= af
}

func (p psarIndicator) Prepare(current *entity.BinCandle, oldData []entity.BinCandle, opts *entity.ChartOption) {
	previous := oldData[len(oldData)-1]
	if previous.Psar == 0.0 {
		if current.Low < previous.Low {
			opts.Low = current.Low
		} else {
			opts.Low = current.Low
		}
		if current.High > previous.High {
			opts.High = current.High
		} else {
			opts.High = current.High
		}

		if current.Close > previous.Close {
			current.Psar = opts.Low
		} else {
			current.Psar = opts.High
		}
		opts.AF = 0.02
	} else {
		if previous.HAClose < previous.Psar {
			sar := previous.Psar - opts.AF*(previous.Psar-opts.Low)
			// if current high > sar, set prev low is sar
			if current.High > sar {
				current.Psar = opts.Low
				opts.High = current.High
				opts.AF = 0.02
				return
			}
			current.Psar = sar
			// if new high, update EP, AF and high
			if current.Low < opts.Low {
				opts.Low = current.Low
				if opts.AF < MAX_AF {
					opts.AF += 0.02
				}
			}
			// uptrend
		} else {
			sar := previous.Psar + opts.AF*(opts.High-previous.Psar)
			// if current low < sar, set prev high is sar
			if current.Low < sar {
				current.Psar = opts.High
				opts.Low = current.Low
				opts.AF = 0.02
				return
			}
			current.Psar = sar
			// if new low, update EP, AF and low
			if current.High > opts.High {
				opts.High = current.High
				if opts.AF < MAX_AF {
					opts.AF += 0.02
				}
			}
		}
	}
}

func (p psarIndicator) CalculateCurrent(candle *entity.BinCandle, opt *entity.ChartOption) {

}

func CalculateHaPsarSingle(current, previous *entity.BinCandle, opts *entity.ChartOption) {

}
