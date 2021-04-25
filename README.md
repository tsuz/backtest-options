# backtest-options

Backtest options strategy using historical data from livevol

[![Build Status](https://travis-ci.org/tsuz/backtest-options.svg?branch=main)](https://travis-ci.org/tsuz/backtest-options) [![codecov](https://codecov.io/gh/tsuz/backtest-options/branch/main/graph/badge.svg?token=pYlB8aWJzB)](https://codecov.io/gh/tsuz/backtest-options)


## Examples

To import data from live vol into your system,

```
go build
> ./backtest-options import /<livevol-dir>

Successfully imported

```

to use the data for strategy, run

```
> ./backtest-options strategy pip  --minCallDTE=4 --minPutDTE=150 

+------------+------------+------------------+------------------+--------------+----------------------+-------------+--------------+------------+---------------+----------------+-------------------+
| OPEN DATE  | CLOSE DATE |   CALL PRODUCT   |   PUT PRODUCT    | TOTAL PROFIT | COVERED CALL PREMIUM | PUT OPEN PX | PUT CLOSE PX | PUT PROFIT | STOCK OPEN PX | STOCK CLOSE PX | CUMULATIVE PROFIT |
+------------+------------+------------------+------------------+--------------+----------------------+-------------+--------------+------------+---------------+----------------+-------------------+
| 2005-01-10 | 2005-01-22 | 121 C 2005-01-22 | 113 P 2005-06-18 |       -154.5 |                 0.25 |        2.08 |         2.55 |       0.48 |       118.845 |        116.575 |            -154.5 |
| 2005-01-24 | 2005-02-19 | 119 C 2005-02-19 | 111 P 2005-09-17 |        154.5 |                0.625 |        3.03 |         1.78 |      -1.25 |       116.575 |        118.745 |                 0 |
| 2005-02-22 | 2005-03-19 | 121 C 2005-03-19 | 113 P 2005-09-17 |           19 |                0.475 |        2.23 |         2.35 |       0.13 |       118.745 |        118.335 |                19 |

....

| 2008-01-22 | 2008-02-16 | 134 C 2008-02-16 | 125 P 2008-06-21 |          357 |                2.885 |        5.88 |         4.03 |      -1.85 |       131.465 |            134 |              1840 |
| 2008-02-19 | 2008-03-22 | 138 C 2008-03-22 | 128 P 2008-09-20 |        167.5 |                2.485 |        6.83 |         6.43 |      -0.40 |       135.225 |        134.815 |            2007.5 |
| 2008-03-24 | 2008-03-31 | 138 C 2008-03-31 | 128 P 2008-09-20 |         -151 |                 0.44 |        6.43 |         6.93 |       0.50 |       134.815 |        132.365 |            1856.5 |
+------------+------------+------------------+------------------+--------------+----------------------+-------------+--------------+------------+---------------+----------------+-------------------+
+-------------------+------------------+--------------+-------------------+
|   TOTAL PROFIT    | TOTAL EXECUTIONS | MAX DRAWDOWN |    BUY & HOLD     |
+-------------------+------------------+--------------+-------------------+
| 1856.50 (15.62 %) |               46 |      1105.26 | 1352.00 (11.38 %) |
+-------------------+------------------+--------------+-------------------+
```


## Strategies

Strategies are run like below

- [x] Covered Call
- [x] PIP, an index hedging strategy (near term covered call + far out long put)
- [ ] The Wheel (short put -> assignment -> short call -> taken away -> short put)
- [ ] Iron Condor at high IV
- [ ] Short Straddles at high IV
- [ ] Long Put at low IV

## Features

- [x] Outputs meta data of the strategy with cumulative profit
- [x] Outputs each execution row as detail
- [ ] Add IV rank to past data and support IV rank in strategies parameter
- [ ] Add Graphs for visual representation
- [ ] Improve backtest performance
- [ ] Ability to export as CSV
- [ ] Ability to run multiple strategies side by side for a comparison
- [ ] Import multiple symbols from CBOE
- [ ] Import Nikkei 225 Options data


## Strategies parameter

### Covered Call

TODO

### PIP Strategy

| Param | Comment | Default |
|--|--|--|
| minPutDTE | MinPutExpDTE is the minimum number of DTE until the next expiry for the put option | 150 |
| minCallDTE | MinCallExpDTE is the minimum number of DTE until the next expiry for the call option | 4 |
