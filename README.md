<h1 align="center">
  <img width="150" src="https://i.ibb.co/3Y1sJDQ/c175b5cd-93cd-4942-b040-d4f4b5abd2b2.jpg" alt="logo" />
  <p>Voltra</p>
</h1>
<p align="center">A cross-market cryptocurrency volatility trading bot written in go.</p>

> [!NOTE]  
> This project has reached its alpha phase. You are invited to report and/or solve bugs!

## Strategy
The bot analyses changes in price of all coins across all supported marketplaces and places trades on the most volatile ones.

## Features
- [x] Supports multiple cryptocurrency marketplaces
- [x] Highly [configurable](./config.example.yml)
    - Configurable stop loss (SL) and take profit (TP)
    - Trailing SL and TP
    - Limit amount of coins to trade
    - Option to include or exclude specific pairs
    - Option to reinvest profits
    - Configurable intervals, trading fees etc.
    - ...
- [x] Optimized for speed and efficiency
    - Selling and buying coins happens in different goroutines (= lightweight threads)
    - Marketplace API requests are reduced to a minimum
- [x] Trade history is stored to a (local) database. 
- [x] Production-grade logging
- [ ] Receive status updates via Telegram or Discord **- coming very soon!**
- [ ] [Request a feature!](https://github.com/sleeyax/voltra/issues/new?assignees=&labels=feature&projects=&template=feature_request.md&title=)

## Supported markets
- [x] Binance
- [ ] [Add marketplace](#add-marketplace)

### Add marketplace
You can request support for a new marketplace by [creating an issue](https://github.com/sleeyax/voltra/issues/new?assignees=&labels=feature&projects=&template=feature_request.md&title=).

If you're a developer, you can add support for a new marketplace by implementing the `Market` interface [here](https://github.com/sleeyax/voltra/blob/main/internal/market/market.go).
See the [Binance](https://github.com/sleeyax/gvoltra/blob/main/internal/market/binance.go) implementation as an example. Comment on the relevant issue if you need help.

## Disclaimer
Trading cryptocurrencies carries risks, including the potential for loss. 
This trading bot is a tool in your toolbox and should not be treated as financial advice. 
Past performance doesn't guarantee future results; markets are volatile, and technical issues or regulatory changes may affect outcomes. 
Users should trade responsibly, understand the risks, and use the bot at their own discretion. 
None of our contributors are liable for any losses incurred.

## Credits
Inspired by [CyberPunkMetalHead/Binance-volatility-trading-bot](https://github.com/CyberPunkMetalHead/Binance-volatility-trading-bot) and [its many forks](https://useful-forks.github.io/?repo=CyberPunkMetalHead/Binance-volatility-trading-bot).
