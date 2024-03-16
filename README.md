<h1 align="center">
  <img width="150" src="https://i.ibb.co/3Y1sJDQ/c175b5cd-93cd-4942-b040-d4f4b5abd2b2.jpg" alt="logo" />
  <p>Voltra</p>
</h1>
<p align="center">A cross-market cryptocurrency volatility trading bot written in go.</p>

> [!WARNING]  
> Trading cryptocurrency carries risks. This bot is a tool in your toolbox, not financial advice. Use at your own discretion.

> [!NOTE]  
> This project has entered its alpha phase. We welcome you to report and/or help resolve any bugs!

![build status](https://github.com/sleeyax/voltra/actions/workflows/build_and_test.yml/badge.svg)

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

## Getting started

### Docker
You can run the bot on any platform or cloud provider that supports Docker.

First, pull the latest image from GitHub's container registry as follows:
```sh
$ docker pull ghcr.io/sleeyax/voltra:latest
```

<details>
  <summary>OR click here for instructions to build from source</summary>
  Clone the source code and build the docker image locally:

  ```sh
  $ git clone https://github.com/sleeyax/voltra.git
  $ cd voltra
  $ docker build --tag ghcr.io/sleeyax/voltra:latest .
  ```
</details>

Then, to run the bot you have the option to run with our without data persistence. When you opt for the latter, any outputted data such as your database will be deleted as soon as the container is removed.

- With full data persistence:

```sh
$ docker run --name voltra --volume ./config.yml:/bot/config.yml:ro --volume ./data:/bot/data -it sleeyax/voltra:latest
```

Alternatively, you can store your config file in the `data` directory and only mount that directory:

```sh
$ docker run --name voltra --volume ./data:/bot/data -it sleeyax/voltra:latest
```

- Without data persistence:

```sh
$ docker run --name voltra --volume ./config.yml:/bot/config.yml:ro -it sleeyax/voltra:latest
```

## Credits
Inspired by [CyberPunkMetalHead/Binance-volatility-trading-bot](https://github.com/CyberPunkMetalHead/Binance-volatility-trading-bot) and [its many forks](https://useful-forks.github.io/?repo=CyberPunkMetalHead/Binance-volatility-trading-bot).
