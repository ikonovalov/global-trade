## Cryptocurrency exchange crafted client. Hold, Buy, Sell.

### Build
```bash
make
```

```
> ./gtr --help
usage: yobit [<flags>] <command> [<args> ...]

Yobit cryptocurrency exchange crafted client.

Flags:
  --help     Show context-sensitive help (also try --help-long and --help-man).
  --version  Show application version.
  --verbose  Print additional information

Commands:
  help [<command>...]
    Show help.

  init <secret> <key>
    Initialize nonce and keys container

  markets [<cryptocurrency>]
    (m) Show all listed tickers on the Yobit

  ticker [<pairs>]
    (tc) Command provides statistic data for the last 24 hours.

  depth [<pairs>] [<limit>]
    (d) Command returns information about lists of active orders for selected pairs.

  trades [<pairs>] [<limit>]
    (tr) Command returns information about the last transactions of selected pairs.

  wallets [<base-currency>]
    (w) Command returns information about user's balances and privileges of API-key as well as server time.

  active-orders <pair>
    (ao) Show active orders

  trade-history <pair>
    (th) Trade history

  buy <pair> <rate> <amount>
    (b) Buy on stock exchange

  sell <pair> <rate> <amount>
    (s) Sell on stock exchange

  cancel <order_id>
    (c) Cancels the chosen order
```
MIT License