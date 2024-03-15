package config

type LogLevel string

const (
	DebugLevel  LogLevel = "debug"
	InfoLevel   LogLevel = "info"
	WarnLevel   LogLevel = "warning"
	ErrorLevel  LogLevel = "error"
	SilentLevel LogLevel = "silent"
)

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

	// The minimum log level.
	// Recommended to set this to `info` in production and `debug` in development.
	LogLevel LogLevel `mapstructure:"log_level"`

	// The minimum database log level.
	// Recommended to set this to `silent` in all environments and only set this to `debug` to log all executed SQL statements in development when necessary.
	// Defaults to LogLevel if not set.
	DatabaseLogLevel LogLevel `mapstructure:"database_log_level"`
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

	// Allows the bot to dynamically adjust the trade Quantity based on the profit/loss of all trades during the current session.
	EnableDynamicQuantity bool `mapstructure:"enable_dynamic_quantity"`

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

	// Trading fee for the maker in % per trade.
	// When you place an order that goes on the order book partially or fully, such as a limit order, any subsequent trades coming from that order will be maker trades.
	// These orders add volume to the order book, help to make the market, and are therefore termed makers for any subsequent trades.
	//
	// Binance:
	//  - If using BNB for fees, set this to 0.075 and make sure have enough BNB in your account.
	TradingFeeMaker float64 `mapstructure:"trading_fee_maker"`

	// Trading fee for the taker in % per trade.
	// When you place an order that trades immediately before going on the order book, you are a taker.
	// This is regardless of whether you partially or fully fulfill the order.
	// Trades from market orders are always takers, as market orders never go on the order book.
	// These trades are "taking" volume off the order book, and therefore are taker trades.
	//
	// Binance:
	//  - If using BNB for fees, set this to 0.075 and make sure have enough BNB in your account.
	TradingFeeTaker float64 `mapstructure:"trading_fee_taker"`

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
