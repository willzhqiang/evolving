package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"ths/ths"
)

var (
	symbolPattern = regexp.MustCompile(`^\d{6}$`)
	pricePattern  = regexp.MustCompile(`^(None|\d+(\.\d{1,3})?)$`)
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	client := ths.GetThsClient()
	simClient := ths.GetThsSimClient()
	switch os.Args[1] {
	case "diag":
		runDiag(client)
	case "account":
		out, err := client.GetAccountInfo()
		printResult("account", out, err)
	case "holdings":
		fs := flag.NewFlagSet("holdings", flag.ExitOnError)
		asset := fs.String("asset", "stock", "asset type: stock, sciTech, gem")
		mustParse(fs, os.Args[2:])
		validateAsset(*asset)
		out, err := client.GetHoldingShares(*asset)
		printResult("holdings", out, err)
	case "entrust":
		fs := flag.NewFlagSet("entrust", flag.ExitOnError)
		asset := fs.String("asset", "stock", "asset type: stock, sciTech, gem")
		revocable := fs.Bool("revocable", true, "show revocable entrusts")
		mustParse(fs, os.Args[2:])
		validateAsset(*asset)
		out, err := client.GetEntrust(*asset, *revocable)
		printResult("entrust", out, err)
	case "buy":
		fs := flag.NewFlagSet("buy", flag.ExitOnError)
		symbol := fs.String("symbol", "", "six-digit stock code")
		qty := fs.Int("qty", 0, "quantity")
		price := fs.String("price", "", `limit price or "None"`)
		asset := fs.String("asset", "", "asset type; auto-detect when empty")
		yes := fs.Bool("yes-live-trade", false, "required for live buy")
		mustParse(fs, os.Args[2:])
		requireLive(*yes, "buy")
		validateOrder(*symbol, *qty, *price, *asset)
		out, err := client.Buy(*symbol, *qty, *price, *asset)
		printResult("buy", out, err)
	case "sell":
		fs := flag.NewFlagSet("sell", flag.ExitOnError)
		symbol := fs.String("symbol", "", "six-digit stock code")
		qty := fs.Int("qty", 0, "quantity")
		price := fs.String("price", "", `limit price or "None"`)
		asset := fs.String("asset", "", "asset type; auto-detect when empty")
		yes := fs.Bool("yes-live-trade", false, "required for live sell")
		mustParse(fs, os.Args[2:])
		requireLive(*yes, "sell")
		validateOrder(*symbol, *qty, *price, *asset)
		out, err := client.Sell(*symbol, *qty, *price, *asset)
		printResult("sell", out, err)
	case "revoke":
		fs := flag.NewFlagSet("revoke", flag.ExitOnError)
		asset := fs.String("asset", "stock", "asset type: stock, sciTech, gem")
		all := fs.Bool("all", false, "allow all revocable orders for the asset to be revoked")
		yes := fs.Bool("yes-live-trade", false, "required for live revoke")
		mustParse(fs, os.Args[2:])
		requireLive(*yes, "revoke")
		if !*all {
			fail("production revoke is currently all-or-nothing; pass --all only after confirming no unrelated revocable orders exist")
		}
		validateAsset(*asset)
		out, err := client.RevokeEntrust(*asset)
		printResult("revoke", out, err)
	case "sim-account":
		out, err := simClient.GetAccountInfo()
		printResult("sim_account", out, err)
	case "sim-holdings":
		fs := flag.NewFlagSet("sim-holdings", flag.ExitOnError)
		asset := fs.String("asset", "stock", "asset type: stock")
		mustParse(fs, os.Args[2:])
		validateSimAsset(*asset)
		out, err := simClient.GetHoldingShares(*asset)
		printResult("sim_holdings", out, err)
	case "sim-entrust":
		fs := flag.NewFlagSet("sim-entrust", flag.ExitOnError)
		asset := fs.String("asset", "stock", "asset type: stock")
		dateRange := fs.String("range", "today", "date range: today, thisWeek, thisMonth, thisSeason, thisYear")
		revocable := fs.Bool("revocable", true, "show revocable entrusts")
		mustParse(fs, os.Args[2:])
		validateSimAsset(*asset)
		validateDateRange(*dateRange)
		out, err := simClient.GetEntrust(*asset, *dateRange, *revocable)
		printResult("sim_entrust", out, err)
	case "sim-buy":
		fs := flag.NewFlagSet("sim-buy", flag.ExitOnError)
		symbol := fs.String("symbol", "", "six-digit stock code")
		qty := fs.Int("qty", 0, "quantity")
		price := fs.String("price", "", `limit price or "None"`)
		asset := fs.String("asset", "stock", "asset type: stock")
		yes := fs.Bool("yes-sim-trade", false, "required for simulated buy")
		mustParse(fs, os.Args[2:])
		requireSim(*yes, "sim-buy")
		validateOrder(*symbol, *qty, *price, *asset)
		validateSimAsset(*asset)
		out, err := simClient.IssuingEntrust("buy", *asset, *symbol, *price, *qty)
		printResult("sim_buy", out, err)
	case "sim-sell":
		fs := flag.NewFlagSet("sim-sell", flag.ExitOnError)
		symbol := fs.String("symbol", "", "six-digit stock code")
		qty := fs.Int("qty", 0, "quantity")
		price := fs.String("price", "", `limit price or "None"`)
		asset := fs.String("asset", "stock", "asset type: stock")
		yes := fs.Bool("yes-sim-trade", false, "required for simulated sell")
		mustParse(fs, os.Args[2:])
		requireSim(*yes, "sim-sell")
		validateOrder(*symbol, *qty, *price, *asset)
		validateSimAsset(*asset)
		out, err := simClient.IssuingEntrust("sell", *asset, *symbol, *price, *qty)
		printResult("sim_sell", out, err)
	case "sim-revoke":
		fs := flag.NewFlagSet("sim-revoke", flag.ExitOnError)
		asset := fs.String("asset", "stock", "asset type: stock")
		contractNo := fs.String("contract-no", "", "contract number to revoke")
		all := fs.Bool("all", false, "revoke all simulated buy/sell orders")
		yes := fs.Bool("yes-sim-trade", false, "required for simulated revoke")
		mustParse(fs, os.Args[2:])
		requireSim(*yes, "sim-revoke")
		validateSimAsset(*asset)
		revokeType := "contractNo"
		if *all {
			revokeType = "allBuyAndSell"
		} else if *contractNo == "" {
			fail("sim-revoke requires --contract-no or --all")
		}
		out, err := simClient.RevokeEntrust(revokeType, *asset, *contractNo)
		printResult("sim_revoke", out, err)
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage:
  evolving-ths diag
  evolving-ths account
  evolving-ths holdings [-asset stock|sciTech|gem]
  evolving-ths entrust [-asset stock|sciTech|gem] [-revocable=true|false]
  evolving-ths buy --symbol 600000 --qty 100 --price 1.23 [--asset stock] --yes-live-trade
  evolving-ths sell --symbol 600000 --qty 100 --price 99.99 [--asset stock] --yes-live-trade
  evolving-ths revoke --asset stock --all --yes-live-trade
  evolving-ths sim-account
  evolving-ths sim-holdings [-asset stock]
  evolving-ths sim-entrust [-asset stock] [-range today] [-revocable=true|false]
  evolving-ths sim-buy --symbol 600000 --qty 100 --price 1.23 --yes-sim-trade
  evolving-ths sim-sell --symbol 600000 --qty 100 --price 99.99 --yes-sim-trade
  evolving-ths sim-revoke --contract-no 123456 --yes-sim-trade
  evolving-ths sim-revoke --all --yes-sim-trade

Live commands operate the currently logged-in broker UI. Confirm the visible app
state before running them. Production revoke is all-or-nothing in the upstream
AppleScript, so the CLI requires --all as an explicit acknowledgement.
`)
}

func runDiag(client interface {
	IsClientLoggedIn() (bool, error)
	IsBrokerLoggedIn() (bool, error)
}) {
	clientLoggedIn, err := client.IsClientLoggedIn()
	fmt.Printf("client_logged_in=%v err=%v\n", clientLoggedIn, err)
	brokerLoggedIn, err := client.IsBrokerLoggedIn()
	fmt.Printf("broker_logged_in=%v err=%v\n", brokerLoggedIn, err)
	if err != nil {
		os.Exit(1)
	}
}

func printResult(label string, out string, err error) {
	fmt.Printf("%s_err=%v\n", label, err)
	fmt.Printf("%s_out=%s\n", label, out)
	if err != nil {
		os.Exit(1)
	}
}

func mustParse(fs *flag.FlagSet, args []string) {
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}
}

func requireLive(yes bool, action string) {
	if !yes {
		fail("%s is a live broker action; rerun with --yes-live-trade after confirming scope", action)
	}
}

func requireSim(yes bool, action string) {
	if !yes {
		fail("%s mutates the simulated trading account; rerun with --yes-sim-trade after confirming scope", action)
	}
}

func validateOrder(symbol string, qty int, price string, asset string) {
	if !symbolPattern.MatchString(symbol) {
		fail("symbol must be a six-digit code")
	}
	if qty <= 0 || qty%100 != 0 {
		fail("qty must be a positive board lot multiple of 100")
	}
	if !pricePattern.MatchString(price) {
		fail(`price must be a decimal with up to 3 places or "None"`)
	}
	if asset != "" {
		validateAsset(asset)
	}
}

func validateAsset(asset string) {
	switch asset {
	case "stock", "sciTech", "gem":
		return
	default:
		fail("asset must be stock, sciTech, or gem")
	}
}

func validateSimAsset(asset string) {
	if asset != "stock" {
		fail("simulation currently supports only stock")
	}
}

func validateDateRange(dateRange string) {
	switch dateRange {
	case "today", "thisWeek", "thisMonth", "thisSeason", "thisYear":
		return
	default:
		fail("range must be today, thisWeek, thisMonth, thisSeason, or thisYear")
	}
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	os.Exit(2)
}
