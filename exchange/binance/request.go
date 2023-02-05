package binance

import (
	"anidiot.com/auto-trade/entity"
	"context"
	binExchange "github.com/adshao/go-binance/v2"
)

var client *binExchange.Client

func InitBinanceClient() {
	apiKey := "dummyAPIKey"
	secretKey := "dummySecretKey"
	client = binExchange.NewClient(apiKey, secretKey)
}

func getOldData(coinName string) ([]entity.BinCandle, error) {
	klines, err := client.NewKlinesService().Symbol(coinName).
		Interval(interval).Limit(200).Do(context.Background())
	if err != nil {
		return nil, err
	}
	data := make([]entity.BinCandle, 0)
	for i := range klines{
		data = append(data, entity.CastCandle(klines[i]))
	}
	return data, nil
}
