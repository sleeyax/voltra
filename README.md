# Volatility Trading Bot
A cross-market cryptocurrency trading bot written in go.

Inspired by [CyberPunkMetalHead/Binance-volatility-trading-bot](https://github.com/CyberPunkMetalHead/Binance-volatility-trading-bot).

> [!WARNING]  
> This project is still under active development. Use at your own risk.

## Strategy
The bot analyses changes in price of all coins across all supported marketplaces and places trades on the most volatile ones.

## Features
- [x] Supports multiple cryptocurrency marketplaces
- [x] Highly [configurable](./config.yml.example)
    - Trailing SL and TP
    - Limit amount of coins to trade
    - Include or exclude specific pairs
    - Configurable intervals, trading fees etc.
    - ...
- [x] Optimized for speed and efficiency
    - Selling and buying coins happens in different goroutines (= lightweight threads)
    - Marketplace API requests are reduced to a minimum
- [x] Sessions are persisted to a (local) database
- [x] Production-grade logging
- [ ] Receive status updates via Telegram or Discord
- [ ] [Request a feature!](https://github.com/sleeyax/go-crypto-volatility-trading-bot/issues/new)

## Supported markets
- [x] Binance
- [ ] [Add marketplace](#add-marketplace)

### Add marketplace
You can request a new marketplace to be supported by [creating an issue](https://github.com/sleeyax/go-crypto-volatility-trading-bot/issues/new) on this repository. 

If you're a developer, you can easily add support for a new marketplace by implementing the `Market` interface [here](https://github.com/sleeyax/go-crypto-volatility-trading-bot/blob/main/internal/market/market.go).
See the [Binance](https://github.com/sleeyax/go-crypto-volatility-trading-bot/blob/main/internal/market/binance.go) implementation as an example. Comment on the relevant issue if you need help.
