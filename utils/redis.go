package utils

import (
	"encoding/json"
	"fmt"
	"time"

	"anidiot.com/auto-trade/entity"
	"anidiot.com/auto-trade/log"
	"go.uber.org/zap"

	// "log"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

var rdb *redis.Client
var keyVersion string
var format string

type DB struct{}

func InitRedis(version string) DB {
	keyVersion = version
	rdb = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	format = viper.GetString("entity.time_format")
	return DB{}
}

func (DB) Set1hCandle(candle entity.CoinTicker) {
	timeString := viper.GetString("setting.time")
	key := fmt.Sprint(candle.Name, ".", timeString, ".1h.", keyVersion)
	tickerSet := make([]redis.Z, 0)
	for i, ticker := range candle.H1Chart {
		data, err := json.Marshal(ticker)
		if err == nil {
			tickerSet = append(tickerSet, redis.Z{
				Score:  float64(i),
				Member: data,
			})
		} else {
			log.Logger.Error("error when marshal", zap.Any("ticker", ticker), zap.Error(err))
		}

	}
	rdb.ZAdd(key, tickerSet...)
}

func (DB) Set1hCandleWithTime(candle entity.CoinTicker, timeString string) {
	key := fmt.Sprint(candle.Name, ".", timeString, ".1h.", keyVersion)
	tickerSet := make([]redis.Z, 0)
	for i, ticker := range candle.H1Chart {
		data, err := json.Marshal(ticker)
		if err == nil {
			tickerSet = append(tickerSet, redis.Z{
				Score:  float64(i),
				Member: data,
			})
		} else {
			log.Logger.Error("error when marshal", zap.Any("ticker", ticker), zap.Error(err))
		}

	}
	rdb.ZAdd(key, tickerSet...)
}

func (DB) SetLast1hCandle(name string, candle entity.Candle, score float64) {
	log.Logger.Info("set last candle", zap.String("name", name), zap.Float64("score", score), zap.Any("candle", candle))
	timeString := viper.GetString("setting.time")
	key := fmt.Sprint(name, ".", timeString, ".1h.", keyVersion)
	log.Logger.Info("add redis", zap.String("key", key))
	data, err := json.Marshal(candle)
	if err == nil {
		result := rdb.ZAdd(key, redis.Z{
			Score:  score,
			Member: data,
		})
		log.Println("result", name, result.String())
	} else {
		log.Logger.Error("error when marshal", zap.Any("candle", candle), zap.Error(err))
	}
}

func (DB) SetCandle(key string, candles []entity.Candle, start int64) {
	key = fmt.Sprint(key, keyVersion)
	tickerSet := make([]redis.Z, 0)
	for i, candle := range candles {
		data, err := json.Marshal(candle)
		if err == nil {
			tickerSet = append(tickerSet, redis.Z{
				Score:  float64(start) + float64(i),
				Member: data,
			})
		} else {
			log.Logger.Error("error when marshal", zap.Any("candle", candle), zap.Error(err))
		}
	}
	rdb.ZAdd(key, tickerSet...)
}

func Get1hCandle(name string, num int64) []entity.Candle {
	quan := num - 1
	timeString := viper.GetString("setting.time")
	key := fmt.Sprint(name, ".", timeString, ".1h.", keyVersion)
	resp := rdb.ZRevRange(key, 0, quan)
	result := ParseCmdString(resp.Val())
	if diff := int64(len(result)); diff < num {
		quan = quan - diff
		lastDay := time.Now().AddDate(0, 0, -1).Format(format)
		key := fmt.Sprint(name, ".", lastDay, ".1h.", keyVersion)
		resp := rdb.ZRevRange(key, 0, quan)
		oldResult := ParseCmdString(resp.Val())
		result = append(result, oldResult...)
	}
	return result
}

func ParseCmdString(value []string) []entity.Candle {
	result := make([]entity.Candle, 0)
	for _, strVal := range value {
		var dest entity.Candle
		err := json.Unmarshal([]byte(strVal), &dest)
		if err != nil {
			log.Logger.Error("ParseCmdString issues when decode redis result to candle", zap.String("strVal", strVal), zap.Error(err))
		}
		result = append(result, dest)
	}
	return result
}

func Get1hCandleString(key string, num int64) []string {
	resp := rdb.ZRevRange(key, 0, num)
	return resp.Val()
}

func GetHighestScore(name string) (*entity.Candle, int64) {
	dayString := fmt.Sprint(name, ".", time.Now().Format(format), ".1h.", keyVersion)
	preDayString := fmt.Sprint(name, ".", time.Now().AddDate(0, 0, -1).Format(format), ".1h.", keyVersion)
	checkKey := rdb.Exists(dayString, preDayString)
	switch checkKey.Val() {
	case 2:
		resp := rdb.ZRevRangeWithScores(dayString, 0, 0)
		last := resp.Val()[0]
		var dest entity.Candle
		err := json.Unmarshal(last.Member.([]byte), &dest)
		if err != nil {
			log.Logger.Error("GetHighestScore 2 issues when decode redis result to candle", zap.Any("strVal", resp.Val()[0]), zap.Error(err))
		}
		return &dest, int64(last.Score)
	case 1:
		resp := rdb.ZRevRangeWithScores(preDayString, 0, 0)
		last := resp.Val()[0]
		var dest entity.Candle
		err := json.Unmarshal(last.Member.([]byte), &dest)
		if err != nil {
			log.Logger.Error("GetHighestScore 1 issues when decode redis result to candle", zap.Any("strVal", resp.Val()[0]), zap.Error(err))
		}
		return &dest, int64(last.Score)
	default:
		return nil, 0
	}
	// return nil, 0
}

func (DB) Set15mCandle(name string, candle []entity.ICandle) {
}

func (DB) Get15mCandle(name string, num int64) []byte {
	return nil
}
