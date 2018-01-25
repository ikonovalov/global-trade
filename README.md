## Yobit cryptocurrency exchange crafted client.
```bash
> ./yobit --help
usage: yobit [<flags>] <command> [<args> ...]

Yobit cryptocurrency exchange crafted client.

Flags:
  --help     Show context-sensitive help (also try --help-long and --help-man).
  --version  Show application version.

Commands:
  help [<command>...]
    Show help.

  markets [<cryptocurrency>]
    (m) Show all listed tickers on the Yobit

  ticker [<pairs>]
    (tc) Command provides statistic data for the last 24 hours.

  depth [<pairs>] [<limit>]
    (d) Command returns information about lists of active orders for selected pairs.

  trades [<pairs>] [<limit>]
    (tr) Command returns information about the last transactions of selected pairs.

  wallets [<currency>]
    (w) Command returns information about user's balances and priviledges of API-key as well as server time.

  active-orders <pair>
    (ao) Show active orders

  trade-history <pair>
    (th) Trade history

  trade <pair> <type> <rate> <amount>
    (t) Creating new orders for stock exchange trading

  buy <pair> <rate> <amount>
    (b) Buy on stock exchange

  sell <pair> <rate> <amount>
    (s) Sell on stock exchange

  cancel <order_id>
    (c) Cancells the chosen order

```
MIT License