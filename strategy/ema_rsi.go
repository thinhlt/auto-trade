package strategy

import "anidiot.com/auto-trade/entity"

type RsiEmaStrategy struct {}

func (r RsiEmaStrategy) Run(candle entity.Candle, option *entity.ChartOption) string {
  return ""
}

