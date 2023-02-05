package log

import (
	"fmt"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"sync"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var f *os.File
var mutex = &sync.Mutex{}
var Logger *zap.Logger

func LogFile(fileName ...string) {
	var name string
	mutex.Lock()
	defer mutex.Unlock()
	if len(fileName) == 0 {
		timeName := viper.GetString("setting.time")
		name = fmt.Sprintf("./z/%s.log", timeName)
	} else {
		name = fileName[0]
	}
	if f != nil {
		f.Close()
	}

	var err error
	f, err = os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		log.Fatal(err)
	}
}

func InitLog(filename string) {
	// f, err := os.OpenFile("./z/demo.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	// if err != nil {
	// 	panic("failed to create temporary file")
	// }

	// config := zap.NewProductionConfig()
	// config.OutputPaths = append(config.OutputPaths, "./demo.log")
	// core := zapcore.NewCore(
	// 	zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
	// 	f,
	// 	zap.InfoLevel)
	// Logger = zap.New(core)
	// zap.NewExample()

	//w := zapcore.AddSync(&lumberjack.Logger{
	//	Filename:   filename,
	//	MaxSize:    500, // megabytes
	//	MaxBackups: 5,
	//	MaxAge:     28, // days
	//})
	//cfg := zap.NewProductionEncoderConfig()
	//cfg.EncodeTime = zapcore.RFC3339TimeEncoder
	//core := zapcore.NewCore(
	//	// cfg.OutputPaths = []string{"judger.log"}
	//	zapcore.NewJSONEncoder(cfg),
	//	w,
	//	zap.InfoLevel,
	//)
	//Logger = zap.New(core)
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	Logger, _ = cfg.Build()
	defer Logger.Sync()
}

func Println(value ...interface{}) {
	fmt.Fprintln(f, value...)
}

func Print(value ...interface{}) {
	fmt.Fprint(f, value...)
}

func Printf(format string, value ...interface{}) {
	fmt.Fprintf(f, format, value...)
}
