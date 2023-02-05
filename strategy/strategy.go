package strategy

import (
  "anidiot.com/auto-trade/entity"
)

type Strategy interface {
  Run(candle entity.BinCandle, option *entity.ChartOption) string
}

var strategyList []Strategy

func InitStrategyList() {

}

func RunStrategies(name string, candle entity.BinCandle, option *entity.ChartOption) []string {
  for i := range strategyList{
    strategyList[i].Run(candle, option)
  }
}