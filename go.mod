module anidiot.com/auto-trade

go 1.16

require (
	anidiot.com/common v1.0.0
	github.com/adshao/go-binance/v2 v2.3.3
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible
	github.com/onsi/gomega v1.12.0 // indirect
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
	github.com/syndtr/goleveldb v1.0.0
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
	go.uber.org/zap v1.16.0
	gorm.io/driver/mysql v1.1.2
	gorm.io/gorm v1.21.14
)

replace anidiot.com/common => /mnt/FE945C1FEF00EE48/project/golang/common
