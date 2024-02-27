package config

type Binance struct {
	ApiKey    string `mapstructure:"api_key"`
	SecretKey string `mapstructure:"secret_key"`
}

type Markets struct {
	Binance Binance `mapstructure:"binance"`
}

type ScriptOptions struct {
	TestMode                bool    `mapstructure:"test_mode"`
	LogTrades               bool    `mapstructure:"log_trades"`
	LogFile                 string  `mapstructure:"log_file"`
	AmericanUser            bool    `mapstructure:"american_user"`
	Markets                 Markets `mapstructure:"markets"`
	EnableStructuredLogging bool    `mapstructure:"structured_logging"`
}

type TradingOptions struct {
	PairWith            string   `mapstructure:"pair_with"`
	Quantity            float64  `mapstructure:"quantity"`
	AllowList           []string `mapstructure:"allow_list"`
	DenyList            []string `mapstructure:"deny_list"`
	MaxCoins            int      `mapstructure:"max_coins"`
	TimeDifference      int      `mapstructure:"time_difference"`
	RecheckInterval     int      `mapstructure:"recheck_interval"`
	SellTimeout         int      `mapstructure:"sell_timeout"`
	ChangeInPrice       float64  `mapstructure:"change_in_price"`
	StopLoss            float64  `mapstructure:"stop_loss"`
	TakeProfit          float64  `mapstructure:"take_profit"`
	UseTrailingStopLoss bool     `mapstructure:"use_trailing_stop_loss"`
	TrailingStopLoss    float64  `mapstructure:"trailing_stop_loss"`
	TrailingTakeProfit  float64  `mapstructure:"trailing_take_profit"`
	TradingFee          float64  `mapstructure:"trading_fee"`
	SignallingModules   []string `mapstructure:"signalling_modules"`
}

type Configuration struct {
	ScriptOptions  ScriptOptions  `mapstructure:"script_options"`
	TradingOptions TradingOptions `mapstructure:"trading_options"`
}
