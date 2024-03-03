package config

type Configuration struct {
	// Whether to perform fake or real trades.
	// Setting this to false will use REAL funds, use at your own risk!
	EnableTestMode bool `mapstructure:"enable_test_mode"`

	// Configuration for bot logs.
	LoggingOptions LoggingOptions `mapstructure:"logging_options"`

	// Configuration for supported cryptocurrency exchanges.
	Markets Markets `mapstructure:"markets"`

	// Main configuration for the trading strategy.
	TradingOptions TradingOptions `mapstructure:"trading_options"`
}

type LoggingOptions struct {
	// Enable or disable logging entirely.
	//  Recommended to set this to true in production and development.
	//  You should only set this to false for testing.
	Enable bool `mapstructure:"enable"`

	// Set this to true if you want to use the structured logging format.
	// Recommended to set this to true in production and false in development.
	EnableStructuredLogging bool `mapstructure:"enable_structured_logging"`
}

type Markets struct {
	Binance Binance `mapstructure:"binance"`
}

type Binance struct {
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
}

type TradingOptions struct {
	// Base currency to use for trading.
	// Recommended to use USDT for most trading pairs.
	PairWith string `mapstructure:"pair_with"`

	// Total amount per trade.
	// Your base currency balance must be at least `max_coins` * `quantity`.
	// Binance uses a minimum of 10 USDT per trade, add a bit extra to enable selling if the price drops.
	// Recommended to specify no less than 12 USDT.
	Quantity float64 `mapstructure:"quantity"`

	// The maximum number of coins to buy at a time.
	// For example, if this is set to 3 and the bot has bought 3 different coins, it will not buy any more until it manages to sell one or more of them.
	// Your base currency balance must be at least `max_coins` * `quantity`.
	MaxCoins int `mapstructure:"max_coins"`

	// The amount of time in MINUTES to wait to calculate the difference from the current price.
	// Recommended minimum is 1.
	TimeDifference int `mapstructure:"time_difference"`

	// Number of times to check for TP/SL during each `time_difference`.
	// Binance allows a maximum of 1200 requests per minute per IP.
	// Recommended minimum is 1.
	RecheckInterval int `mapstructure:"recheck_interval"`

	// The amount of time in SECONDS to wait between each try to sell your current coin holdings.
	SellTimeout int `mapstructure:"sell_timeout"`

	// The minimum difference in PERCENTAGE between the previous and current price of a coin to identify it as volatile.
	ChangeInPrice float64 `mapstructure:"change_in_price"`

	// Specify in PERCENTAGE how much you are willing to lose on a coin.
	// For example, if you set this to 5, the bot will sell the coin if it drops 5% below the price at which it was bought.
	StopLoss float64 `mapstructure:"stop_loss"`

	// Specify in PERCENTAGE how much you are looking to gain on a coin.
	// For example, if you set this to 5, the bot will sell the coin if it rises 5% above the price at which it was bought.
	TakeProfit float64 `mapstructure:"take_profit"`

	// Trading fee in % per trade.
	//
	// Binance:
	// - If using 0.75% (using BNB for fees) you must have BNB in your account to cover trading fees.
	// - If using BNB for fees, it MUST be enabled in your Binance 'Dashboard' page (checkbox).
	TradingFee float64 `mapstructure:"trading_fee"`

	// The amount of time in MINUTES to wait before buying the same coin again.
	// This is to prevent buying the same coin multiple times in a short period of time.
	// Set to 0 to disable.
	CoolOffDelay int `mapstructure:"cool_off_delay"`

	// Configuration for trailing stop loss.
	TrailingStopOptions TrailingStopOptions `mapstructure:"trailing_stop_options"`

	// List of tickers to include.
	AllowList []string `mapstructure:"allow_list"`

	// List of trading pairs to exclude.
	DenyList []string `mapstructure:"deny_list"`
}

type TrailingStopOptions struct {
	// Whether to enable trailing stop loss.
	// If true, the bot will automatically move the stop loss up as the price of the coin increases to 'lock-in' a profit.
	// Recommended to set this to true.
	Enable bool `mapstructure:"enable"`

	// When `take_profit` is reached, the `stop_loss` is changed to `trailing_stop_loss` PERCENTAGE below `take_profit` hence 'locking in' the profit.
	TrailingStopLoss float64 `mapstructure:"trailing_stop_loss"`

	// When `take_profit` is reached, the `take_profit` is changed to `trailing_take_profit` PERCENTAGE above the current price.
	TrailingTakeProfit float64 `mapstructure:"trailing_take_profit"`
}
