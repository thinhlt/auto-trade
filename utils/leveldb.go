package utils

import (
	"anidiot.com/auto-trade/entity"
	"anidiot.com/auto-trade/log"
	"encoding/json"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
	"go.uber.org/zap"
	"sync"
	"time"
)

var TIME_PARTITION = 900000
var TIME_PARTITION_30 = 1800000

type LevelDB struct {
	db *leveldb.DB
}

type Record struct {
	key   []byte
	value []byte
}

var writeChan chan Record
var stopChan chan struct{}
var mux = sync.Mutex{}
var opts = opt.WriteOptions{Sync: true}

func InitLevelDB(path string) LevelDB {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		log.Logger.Error("Init leveldb error", zap.Error(err))
	}
	writeChan = make(chan Record, 1)
	stopChan = listenWriteChan(db)
	return LevelDB{db}
}

func listenWriteChan(db *leveldb.DB) chan struct{} {
	batch := new(leveldb.Batch)
	ticker := time.NewTicker(time.Minute)
	doneC := make(chan struct{}, 1)
	go func() {
		for {
			select {
			case record := <-writeChan:
				batch.Put(record.key, record.value)
			case <-ticker.C:
				if batch.Len() > 0 {
					err := db.Write(batch, &opts)
					if err != nil {
						log.Logger.Error("Set to db error", zap.Error(err))
					} else {
						batch = new(leveldb.Batch)
					}
				}
			case <-doneC:
				log.Logger.Info("DB writer stop")
				return
			}
		}
	}()
	return doneC
}

func (LevelDB) Set1hCandle(candle entity.CoinTicker) {
}

func (LevelDB) Set1hCandleWithTime(candle entity.CoinTicker, timeString string) {
}

func (LevelDB) SetLast1hCandle(name string, candle entity.Candle, score float64) {
}

func (LevelDB) SetCandle(key string, candles []entity.Candle, start int64) {
}

// TODO: export calculate day score
func (inst LevelDB) Set15mCandle(name string, candle []entity.ICandle) {
	for i := range candle {
		score := int(candle[i].GetTime()) / TIME_PARTITION
		candleByte := candle[i].ToByteArray()
		key := []byte(fmt.Sprintf("%s_%010d", name, score))
		err := inst.db.Put(key, candleByte, nil)
		if err != nil {
			log.Logger.Error("Set to db error", zap.Error(err),
				zap.String("key", string(key)), zap.String("data", string(candleByte)))
		}
	}
}

func (inst LevelDB) SetLast15mCandle(name string, candle entity.ICandle) {
	score := candle.GetScore()
	candleByte := candle.ToByteArray()
	key := []byte(fmt.Sprintf("%s_%010d", name, score))
	err := inst.db.Put(key, candleByte, nil)
	if err != nil {
		log.Logger.Error("Set to db error", zap.Error(err),
			zap.String("key", string(key)), zap.String("data", string(candleByte)))
	}
}

func (LevelDB) GetLatest(name string, candle []entity.ICandle) ([]entity.ICandle, error) {
	return []entity.ICandle{}, nil
}

func (inst LevelDB) Get15mCandle(name string, currentTime int64, num int) ([]entity.ICandle, error) {
	var start, stop []byte
	score := (int(currentTime) / TIME_PARTITION) - 1
	start = []byte(fmt.Sprintf("%s_%010d", name, score-num))
	stop = []byte(fmt.Sprintf("%s_%010d", name, score))
	result := make([]entity.ICandle, 0)

	queryResult := inst.db.NewIterator(&util.Range{Start: start, Limit: stop}, nil)
	if err := queryResult.Error(); err != nil {
		log.Logger.Error("error when query leveldb", zap.Error(err))
		return nil, err
	}
	for queryResult.Next() {
		val := queryResult.Value()
		var tick entity.BinCandle
		json.Unmarshal(val, &tick)
		result = append(result, tick)
	}
	//log.Logger.Info("get old data success", zap.String("startkey", queryResult.Key()))
	return result, nil
}

func (inst LevelDB) Set30mCandle(name string, candle []entity.ICandle) {
	batch := new(leveldb.Batch)
	for i := range candle {
		score := int(candle[i].GetTime()) / TIME_PARTITION_30
		candleByte := candle[i].ToByteArray()
		str := fmt.Sprintf("%s_%010d", name, score)
		key := []byte(str)
		batch.Put(key, candleByte)
	}
	mux.Lock()
	if err := inst.db.Write(batch, &opts); err != nil {
		log.Logger.Error("Set batch to db error", zap.Error(err),
			zap.String("key", string(name)))
	}
	mux.Unlock()
}

func (inst LevelDB) SetLast30mCandle(name string, candle entity.ICandle) {
	score := candle.GetScore()
	candleByte := candle.ToByteArray()
	str := fmt.Sprintf("%s_%010d", name, score)
	key := []byte(str)
	writeChan <- Record{key: key, value: candleByte}
}

func (inst LevelDB) Get30mCandle(name string, currentTime int64, num int) ([]entity.ICandle, error) {
	var start, stop []byte
	score := (int(currentTime) / TIME_PARTITION_30) - 1
	startStr := fmt.Sprintf("%s_%010d", name, score-(2*num))
	stopStr := fmt.Sprintf("%s_%010d", name, score)
	start = []byte(startStr)
	stop = []byte(stopStr)
	result := make([]entity.ICandle, 0)

	queryResult := inst.db.NewIterator(&util.Range{Start: start, Limit: stop}, nil)
	if err := queryResult.Error(); err != nil {
		log.Logger.Error("error when query leveldb", zap.Error(err))
		return nil, err
	}
	for queryResult.Next() {
		val := queryResult.Value()
		var tick entity.BinCandle
		json.Unmarshal(val, &tick)
		result = append(result, tick)
	}
	if len(result) < num {
		log.Logger.Info("db store less than request", zap.String("symbol", name), zap.Int("lenght", len(result)), zap.Any("result", result))
	}
	result = result[len(result)-num:]
	log.Logger.Info("get old data success", zap.String("name", name),
		zap.Int64("startTime", result[0].GetTime()), zap.Int64("endTime", result[num-1].GetTime()))
	return result, nil
}

func (inst LevelDB) GetData() ([]entity.ICandle, error) {
	result := make([]entity.ICandle, 0)

	queryResult := inst.db.NewIterator(nil, nil)
	if err := queryResult.Error(); err != nil {
		log.Logger.Error("error when query leveldb", zap.Error(err))
		return nil, err
	}
	for queryResult.Next() {
		log.Logger.Info("", zap.String("key", string(queryResult.Key())), zap.String("value", string(queryResult.Value())))
		//inst.db.Delete(queryResult.Key(), nil)
	}
	//log.Logger.Info("get old data success", zap.String("startkey", queryResult.Key()))
	return result, nil
}

func (inst LevelDB) Destroy() {
	stopChan <- struct{}{}
}
