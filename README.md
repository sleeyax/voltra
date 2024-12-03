<h1 align="center">
  <img width="150" src="https://i.ibb.co/3Y1sJDQ/c175b5cd-93cd-4942-b040-d4f4b5abd2b2.jpg" alt="logo" />
  <p>Voltra</p>
</h1>
<p align="center">A cross-market cryptocurrency volatility trading bot written in go.</p>

> [!WARNING]  
> Trading cryptocurrency carries risks. This bot is a tool, not financial advice. Use at your own discretion.

![build status](https://github.com/sleeyax/voltra/actions/workflows/run_tests.yml/badge.svg)

## Project status

**NO LONGER MAINTAINED. CONSIDER USING [gocryptotrader](https://github.com/thrasher-corp/gocryptotrader) INSTEAD!**

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
- [ ] Receive status updates via Telegram or Discord
- [ ] [Request a feature!](https://github.com/sleeyax/voltra/issues/new?assignees=&labels=feature&projects=&template=feature_request.md&title=)

## Supported markets
- [x] Binance
- [ ] [Request marketplace](https://github.com/sleeyax/voltra/issues/new?assignees=&labels=feature,marketplace+request&projects=&template=feature_request.md&title=)

If you're a developer, you can add support for a new marketplace by implementing the `Market` interface [here](https://github.com/sleeyax/voltra/blob/main/internal/market/market.go).
See the [Binance](https://github.com/sleeyax/gvoltra/blob/main/internal/market/binance.go) implementation as an example. Comment on the relevant issue if you need help.

## Getting started

First of all, make sure you have a valid config file. You can use the [example config](./config.example.yml) as a starting point:

```sh
$ cp config.example.yml config.yml
```

Open the file in your text editor of choice and edit it according to your preferred trading strategy. Don't forget to at least configure your API keys!

### Binaries
The easiest way to get started is by downloading the latest release for your operating system from the [releases page](https://github.com/sleeyax/voltra/releases). 
Extract the archive, copy your config file to the same directory and finally run the binary in a terminal/command prompt as follows:

```sh
$ ./voltra
```

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
