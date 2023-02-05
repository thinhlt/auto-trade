package msql_client

import (
	"anidiot.com/auto-trade/constraint"
	"anidiot.com/auto-trade/entity"
	logger "anidiot.com/auto-trade/log"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gomLog "gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type MySqlClient struct {
	db *gorm.DB
}

func (m MySqlClient) Set1hCandle(candle entity.CoinTicker) {
	panic("implement me")
}

func (m MySqlClient) Set1hCandleWithTime(candle entity.CoinTicker, timeString string) {
	panic("implement me")
}

func (m MySqlClient) SetLast1hCandle(name string, candle entity.Candle, score float64) {
	panic("implement me")
}

func (m MySqlClient) SetCandle(key string, candles []entity.Candle, start int64) {
	panic("implement me")
}

func (m MySqlClient) Set15mCandle(name string, candles []entity.ICandle) {
	panic("implement me")
}

func (m MySqlClient) Set30mCandle(name string, candles []entity.BinCandle) {
	if !m.db.Migrator().HasTable(name) {
		tick := entity.BinCandle{
			Name: name,
		}
		m.db.Table(name).AutoMigrate(&tick)
	}
	tx := m.db.Table(name).Save(&candles)
	if tx.Error != nil {
		logger.Logger.Info("Error when write to my sql",
			zap.String("name", name), zap.Any("candles", candles), zap.Error(tx.Error))
	}
}

func (m MySqlClient) SetLast15mCandle(name string, candle entity.ICandle) {
	panic("implement me")
}

func (m MySqlClient) SetLast30mCandle(name string, candle entity.BinCandle) {
	tx := m.db.Table(name).Save(&candle)
	if tx.Error != nil {
		logger.Logger.Info("Error when write last candle to my sql",
			zap.String("name", name), zap.Any("candles", candle), zap.Error(tx.Error))
	}
}

func (m MySqlClient) Get15mCandle(name string, currentTime int64, num int) ([]entity.ICandle, error) {
	panic("implement me")
}

func (m MySqlClient) Get30mCandle(name string, currentTime int64, num int) ([]entity.BinCandle, error) {
	result := make([]entity.BinCandle, 0)
	stop := currentTime / constraint.TIME_PARTITION_30
	start := stop - int64(num) - 10
	tx := m.db.Table(name).Where("score > ? and score < ?", start, stop).Find(&result)
	if tx.Error != nil {
		logger.Logger.Error("Error when query old data", zap.String("name", name), zap.Error(tx.Error))
		return result, tx.Error
	}
	if len(result) > num {
		result = result[len(result)-num:]
	}
	return result, nil
}

func (m MySqlClient) Destroy() {
	logger.Logger.Info("Close db connection")
	db, _ := m.db.DB()
	db.Close()
}

func Init() MySqlClient {
	connectUrl := viper.GetString("my_sql.url")
	f, err := os.OpenFile("./db.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	logInst := gomLog.New(log.New(f, "\r\n", log.LstdFlags), gomLog.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  gomLog.Silent,
		IgnoreRecordNotFoundError: false,
		Colorful:                  true,
	})
	db, err := gorm.Open(mysql.Open(connectUrl), &gorm.Config{Logger: logInst})
	if err != nil {
		logger.Logger.Error("Error when connnect to my sql", zap.Error(err))
	}
	return MySqlClient{db: db}
}
