package log

import (
	"anidiot.com/auto-trade/entity"
	"go.uber.org/zap"
)

type LogStore struct{}

func (log LogStore) Set30mCandle(name string, candles []entity.BinCandle) {
	Logger.Info("Set30mCandle", zap.String("candle name", name), zap.Reflect("1h chart", candles))
}

func (log LogStore) SetLast30mCandle(name string, candle entity.BinCandle) {
	Logger.Info("SetLast30mCandle", zap.String("candle name", name), zap.Reflect("1h chart", candle))
}

func (log LogStore) Get30mCandle(name string, currentTime int64, num int) ([]entity.BinCandle, error) {
	return nil, nil
}

func (log LogStore) Destroy() {
	panic("implement me")
}

// dbstore
func (log LogStore) Set1hCandle(candle entity.CoinTicker) {
	Logger.Info("Set1hCandle", zap.String("candle name", candle.Name), zap.Reflect("1h chart", candle.H1Chart))
}
func (log LogStore) Set1hCandleWithTime(candle entity.CoinTicker, timeString string) {
	Logger.Info("Set1hCandleWithTime", zap.String("candle name", candle.Name), zap.Reflect("1h chart", candle.H1Chart), zap.String("time string", timeString))
}
func (log LogStore) SetLast1hCandle(name string, candle entity.Candle, score float64) {
	Logger.Info("SetLast1hCandle", zap.String("candle name", name), zap.Reflect("candle", candle), zap.Float64("score", score))
}
func (log LogStore) SetCandle(key string, candles []entity.Candle, start int64) {
	Logger.Info("SetCandle", zap.String("key", key), zap.Reflect("candle", candles), zap.Int64("start", start))
}
func (log LogStore) Set15mCandle(name string, candles []entity.ICandle) {
	Logger.Info("Set1hCandle", zap.String("candle name", name), zap.Reflect("1h chart", candles))
}
func (log LogStore) SetLast15mCandle(name string, candle entity.ICandle) {
	Logger.Info("Set1hCandle", zap.String("candle name", name), zap.Reflect("1h chart", candle))
}
func (log LogStore) Get15mCandle(name string, currentTime int64, num int) ([]entity.ICandle, error) {
	return nil, nil
}

// message queue
func (log LogStore) SendMessage(messages []string, channelID int64) {
	if len(messages) > 0 {
		Logger.Info("SendMessage", zap.Strings("messages", messages))
	}
}
func (LogStore) Close() {
	Logger.Info("Close")
}
