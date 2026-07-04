# evolving

MacOS automated trading engine for Tonghuashun (同花顺) via AppleScript.

## Features

- Stock trading (A-share, STAR, GEM)
- Bank-broker fund transfers
- Order management
- Holdings and entrustment queries
- IPO subscription
- Simulation trading

## Structure

```
ths/
├── ascmds.go        # osascript executor
├── client.go        # thsClient (production) singleton
├── client_sim.go    # thsSimClient (simulation) singleton
├── consts.go        # production AppleScript constants (34)
└── consts_sim.go    # simulation AppleScript constants (7)
```

## Requirements

- MacOS
- cliclick >= 4.0.1 (`brew install cliclick`)
- Tonghuashun client

The AppleScript snippets try `/opt/homebrew/bin/cliclick` first and then
`/usr/local/bin/cliclick`, covering both Apple Silicon and Intel Homebrew
defaults.

## Usage

```go
import "evolving/ths"

client := ths.GetThsClient()
simClient := ths.GetThsSimClient()
```

## Maintained fork notes

This fork keeps the original AppleScript approach, but makes the runtime usable
from Go again:

- `ths.Run` now supports the repository's legacy constants, which are stored as
  shell commands like `osascript -e '...'`.
- A small CLI is available under `cmd/evolving-ths` for smoke tests and guarded
  trading experiments.
- Live broker-mutating commands require `--yes-live-trade`.
- Simulated-account mutating commands require `--yes-sim-trade`.
- Production revoke now supports `--contract-no`; use live `--all` only after
  checking that there are no unrelated revocable orders.
- Simulation POC creates a low-price buy entrust, reads back the new contract
  number, then revokes it and verifies the revocable list is empty. If
  contract-number revoke does not clear the new simulated order and the
  pre-check list was empty, it falls back to simulated `--all`.
- Live POC is intentionally narrower: it reads account state, submits one buy
  attempt, reads new revocable entrusts, and revokes only newly detected
  contract numbers. Use the explicit `sell` command for sell-path tests.
- The Tonghuashun macOS app can expose multiple windows through Accessibility.
  This fork selects the visible trading window dynamically instead of assuming
  `window 1`.

Examples:

```bash
go run ./cmd/evolving-ths diag
go run ./cmd/evolving-ths account
go run ./cmd/evolving-ths holdings --asset stock
go run ./cmd/evolving-ths entrust --asset stock --revocable=true

# Simulated account
go run ./cmd/evolving-ths sim-account
go run ./cmd/evolving-ths poc --mode sim --symbol 589850 --qty 100 --buy-price 1.800 --asset stock --yes-sim-trade
go run ./cmd/evolving-ths sim-buy --symbol 589850 --qty 100 --price 1.800 --yes-sim-trade
go run ./cmd/evolving-ths sim-entrust --asset stock --range today --revocable=true
go run ./cmd/evolving-ths sim-revoke --all --yes-sim-trade

# Live broker actions. Confirm current visible UI state first.
go run ./cmd/evolving-ths buy --symbol 159949 --qty 100 --price 2.000 --asset stock --yes-live-trade
go run ./cmd/evolving-ths sell --symbol 159949 --qty 100 --price 2.100 --asset stock --yes-live-trade
go run ./cmd/evolving-ths revoke --asset stock --contract-no 123456 --yes-live-trade
go run ./cmd/evolving-ths revoke --asset stock --all --yes-live-trade
```

### Safety model

This project drives the visible Tonghuashun macOS UI. It is not a broker API.
Treat every mutating command as if a human clicked the same buttons:

- Check login state and account state before every live action.
- Read current revocable entrusts before and after every test.
- Prefer simulated-account tests first.
- Prefer live `revoke --contract-no`; use live `revoke --all` only when the
  current revocable list contains only the intended test orders.
- During non-trading days the real broker may reject orders before creating any
  entrust, for example with `证券交易未初始化`.
- Prefer `poc --mode sim` for end-to-end smoke tests. In current testing,
  `589850` parsed correctly in Tonghuashun simulation while `159949` did not
  always trigger the app's internal market-code resolution.
- For scripts that type a security code into Tonghuashun, direct
  Accessibility `set value` may update the visible field without triggering the
  app's internal quote/market parser. Prefer an already-selected symbol or a
  real keyboard input path when building new flows.

### thsClient (production)

```go
// Client login/logout
client.LoginClient(userID, password string) (bool, error)
client.LogoutClient() (bool, error)
client.IsClientLoggedIn() (bool, error)

// Broker login/logout
client.LoginBroker(brokerName, account, password string) (bool, error)
client.LogoutBroker() (bool, error)
client.IsBrokerLoggedIn() (bool, error)

// Trading
client.Buy(stockCode string, amount int, price, assetType string) (string, error)
client.Sell(stockCode string, amount int, price, assetType string) (string, error)
client.Transfer(transferType string, amount int, bankPassword, brokerPassword string) (bool, error)
client.TransferBank2Broker(amount int, bankPassword, brokerPassword string) (bool, error)
client.TransferBroker2Bank(amount int, bankPassword, brokerPassword string) (bool, error)

// Revoke
client.RevokeAllEntrust() (string, error)
client.RevokeAllBuyEntrust() (string, error)
client.RevokeAllSellEntrust() (string, error)
client.RevokeEntrust(assetType string) (string, error)

// Queries
client.GetAccountInfo() (string, error)
client.GetHoldingShares(assetType string) (string, error)  // stock/sciTech/gem
client.GetEntrust(assetType string, isRevocable bool) (string, error)
client.GetClosedDeals(assetType string) (string, error)
client.GetBids(stockCode, assetType string) (string, error)
client.GetTransferRecords(dateRange string) (string, error)
client.GetCapitalDetails(assetType, dateRange string) (string, error)

// IPO
client.OneKeyIPO() (string, error)
client.GetTodayIPO() (string, error)
client.GetIPO(queryType, dateRange string) (string, error)

// Asset type detection
client.GetAssetType(stockCode string) string  // stock/sciTech/gem
```

### thsSimClient (simulation)

```go
simClient.LoginClient(userID, password string) (bool, error)
simClient.LogoutClient() (bool, error)
simClient.IsClientLoggedIn() (bool, error)
simClient.GetAccountInfo() (string, error)
simClient.GetHoldingShares(assetType string) (string, error)
simClient.GetEntrust(assetType, dateRange string, isRevocable bool) (string, error)
simClient.GetClosedDeals(assetType, dateRange string) (string, error)
simClient.GetCapitalDetails(assetType, dateRange string) (string, error)
simClient.IssuingEntrust(tradingAction, assetType, stockCode, price string, amount int) (string, error)
simClient.RevokeEntrust(revokeType, assetType, contractNo string) (string, error)
```

## Permissions

Mac -> System Settings -> Privacy & Security -> Accessibility / Full Disk Access
- [x] Terminal
- [x] osascript

## License

MIT
