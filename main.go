package main

import (
	"fmt"
	"os"

	"anidiot.com/auto-trade/factory"
	"anidiot.com/auto-trade/indicator"
	"anidiot.com/auto-trade/strategy"
	"go.uber.org/zap"

	"anidiot.com/auto-trade/log"

	binExchange "anidiot.com/auto-trade/exchange/binance"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
)

func init() {
	fmt.Println("init()")
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is config.toml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	fmt.Println("read config")
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config") // name of config file (without extension)
	}

	viper.AddConfigPath(".") // adding home directory as first search path
	viper.AutomaticEnv()     // read in environment variables that match
	viper.SetConfigType("toml")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Printf("error when read config file %v \n", err)
	}
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:          "root-command",
	Short:        "auto-trade",
	Long:         "automation bot",
	SilenceUsage: true,
}

var binanceExChange = &cobra.Command{
	Use:   "binance",
	Short: "binance",
	Long:  "binance",
	Run: func(cmd *cobra.Command, args []string) {
		// db.Init()
		// sc.StartScraping(context.Background())
		fmt.Println("PID:", os.Getpid())
		env := viper.GetString("setting.env")
		if env == "production" {
			log.InitLog("auto_bot.log")
			factory.InitProduction()
		} else {
			log.InitLog("auto_bot.log")
			factory.InitDevelopment()
		}
		list := viper.GetStringSlice("binance.watch_list")
		log.Logger.Info("binance", zap.Any("list of binance coin", list))
		binExchange.InitBinanceClient()
		indicator.InitIndicatorList()
		strategy.InitStrategyList()
		mapWatchingCoin := binExchange.InitExchangeBot(list)

		binExchange.Simple(mapWatchingCoin)
		log.Logger.Info("=================================================================================")
	},
}

func Execute() {
	RootCmd.AddCommand(binanceExChange)

	if err := RootCmd.Execute(); err != nil {
		fmt.Println("Start program failed", err)
		os.Exit(-1)
	}
}

func main() {
	Execute()
}

// 1442602480:AAFbZwD5543nz7l2lNYMMbp-L8rNApIBAnw
