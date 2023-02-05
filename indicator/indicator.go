package indicator

import "anidiot.com/auto-trade/entity"

type Indicator interface {
  Calculate(oldData []entity.BinCandle, opt *entity.ChartOption)
  Prepare(candle *entity.BinCandle, oldData []entity.BinCandle, opt *entity.ChartOption)
  CalculateCurrent(candle *entity.BinCandle, opt *entity.ChartOption)
}

type List struct {
  list []Indicator
}

var indicatorIterator []Indicator

func InitIndicatorList() List {
  indicatorIterator = make([]Indicator, 0)
  indicatorIterator = append(indicatorIterator, bbIndicator{})
  indicatorIterator = append(indicatorIterator, emaIndicator{})
  indicatorIterator = append(indicatorIterator, kdjIndicator{})
  indicatorIterator = append(indicatorIterator, rsiIndicator{})
  indicatorIterator = append(indicatorIterator, HeikinAshiCandle{})
  indicatorIterator = append(indicatorIterator, psarIndicator{})
  return List{indicatorIterator}
}

func Calculate(oldData []entity.BinCandle) entity.ChartOption{
  opt := entity.ChartOption{}
  for _ ,indicator := range indicatorIterator {
    indicator.Calculate(oldData, &opt)
  }
  return opt
}

func CalculateNewTicker(name string, candle *entity.BinCandle, oldData []entity.BinCandle, opt *entity.ChartOption){
  for _ ,indicator := range indicatorIterator {
    indicator.Prepare(candle, oldData, opt)
  }
}

func CalculateCurrent(name string, candle *entity.BinCandle, opt *entity.ChartOption){
  for _ ,indicator := range indicatorIterator {
    indicator.CalculateCurrent(candle, opt)
  }
}
