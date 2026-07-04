package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

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
	case "poc":
		fs := flag.NewFlagSet("poc", flag.ExitOnError)
		mode := fs.String("mode", "sim", "poc mode: sim or live")
		symbol := fs.String("symbol", "589850", "six-digit stock code")
		qty := fs.Int("qty", 100, "quantity")
		buyPrice := fs.String("buy-price", "1.800", `buy limit price or "None"`)
		asset := fs.String("asset", "stock", "asset type")
		yesLive := fs.Bool("yes-live-trade", false, "required for live poc")
		yesSim := fs.Bool("yes-sim-trade", false, "required for sim poc")
		mustParse(fs, os.Args[2:])
		validateOrder(*symbol, *qty, *buyPrice, *asset)
		switch *mode {
		case "sim":
			requireSim(*yesSim, "poc --mode sim")
			validateSimAsset(*asset)
			runSimPOC(simClient, *symbol, *qty, *buyPrice, *asset)
		case "live":
			requireLive(*yesLive, "poc --mode live")
			validateAsset(*asset)
			runLivePOC(client, *symbol, *qty, *buyPrice, *asset)
		default:
			fail("mode must be sim or live")
		}
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
		contractNo := fs.String("contract-no", "", "contract number to revoke")
		all := fs.Bool("all", false, "allow all revocable orders for the asset to be revoked")
		yes := fs.Bool("yes-live-trade", false, "required for live revoke")
		mustParse(fs, os.Args[2:])
		requireLive(*yes, "revoke")
		validateAsset(*asset)
		var out string
		var err error
		if *all {
			out, err = client.RevokeEntrust(*asset)
		} else {
			if *contractNo == "" {
				fail("revoke requires --contract-no or --all")
			}
			out, err = client.RevokeEntrustByContractNo(*asset, *contractNo)
		}
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
  evolving-ths poc --mode sim --symbol 589850 --qty 100 --buy-price 1.800 --yes-sim-trade
  evolving-ths poc --mode live --symbol 159949 --qty 100 --buy-price 2.000 --yes-live-trade
  evolving-ths account
  evolving-ths holdings [-asset stock|sciTech|gem]
  evolving-ths entrust [-asset stock|sciTech|gem] [-revocable=true|false]
  evolving-ths buy --symbol 600000 --qty 100 --price 1.23 [--asset stock] --yes-live-trade
  evolving-ths sell --symbol 600000 --qty 100 --price 99.99 [--asset stock] --yes-live-trade
  evolving-ths revoke --asset stock --contract-no 123456 --yes-live-trade
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
AppleScript, but this fork supports contract-number revocation for live broker
orders. Use --all only after checking there are no unrelated revocable orders.
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

type liveClient interface {
	GetAccountInfo() (string, error)
	GetEntrust(assetType string, isRevocable bool) (string, error)
	Buy(stockCode string, amount int, price, assetType string) (string, error)
	RevokeEntrustByContractNo(assetType, contractNo string) (string, error)
}

type simClient interface {
	GetAccountInfo() (string, error)
	GetEntrust(assetType, dateRange string, isRevocable bool) (string, error)
	IssuingEntrust(tradingAction, assetType, stockCode, price string, amount int) (string, error)
	RevokeEntrust(revokeType, assetType, contractNo string) (string, error)
}

func runLivePOC(client liveClient, symbol string, qty int, buyPrice string, asset string) {
	fmt.Println("poc_mode=live")
	accountOut, accountErr := step("live_account", client.GetAccountInfo)
	_, _ = accountOut, accountErr
	beforeOut, beforeErr := step("live_entrust_before", func() (string, error) {
		return client.GetEntrust(asset, true)
	})
	beforeContracts := extractContracts(beforeOut)
	buyOut, buyErr := step("live_buy", func() (string, error) {
		return client.Buy(symbol, qty, buyPrice, asset)
	})
	_, _ = buyOut, buyErr
	if buyErr != nil {
		step("live_entrust_after_buy_failure", func() (string, error) {
			return client.GetEntrust(asset, true)
		})
		fmt.Println("poc_stopped=live-buy-failed")
		return
	}
	afterOut, afterErr := step("live_entrust_after", func() (string, error) {
		return client.GetEntrust(asset, true)
	})
	if beforeErr != nil || afterErr != nil {
		fmt.Println("poc_revoke_skipped=entrust-read-failed")
		return
	}
	newContracts := diffContracts(beforeContracts, extractContracts(afterOut))
	fmt.Printf("poc_new_contracts=%s\n", strings.Join(newContracts, ","))
	for _, contract := range newContracts {
		step("live_revoke_"+contract, func() (string, error) {
			return client.RevokeEntrustByContractNo(asset, contract)
		})
	}
	step("live_entrust_final", func() (string, error) {
		return client.GetEntrust(asset, true)
	})
}

func runSimPOC(client simClient, symbol string, qty int, buyPrice string, asset string) {
	fmt.Println("poc_mode=sim")
	step("sim_account", client.GetAccountInfo)
	beforeOut, beforeErr := step("sim_entrust_before", func() (string, error) {
		return client.GetEntrust(asset, "today", true)
	})
	beforeContracts := extractContracts(beforeOut)
	if beforeErr != nil {
		fmt.Println("poc_stopped=sim-entrust-before-failed")
		return
	}
	buyOut, buyErr := step("sim_buy", func() (string, error) {
		return client.IssuingEntrust("buy", asset, symbol, buyPrice, qty)
	})
	_, _ = buyOut, buyErr
	if buyErr != nil {
		fmt.Println("poc_stopped=sim-buy-failed")
		return
	}
	afterOut, afterErr := step("sim_entrust_after", func() (string, error) {
		return client.GetEntrust(asset, "today", true)
	})
	if afterErr != nil {
		fmt.Println("poc_revoke_skipped=entrust-read-failed")
		return
	}
	newContracts := diffContracts(beforeContracts, extractContracts(afterOut))
	fmt.Printf("poc_new_contracts=%s\n", strings.Join(newContracts, ","))
	if len(newContracts) == 0 {
		fmt.Println("poc_revoke_skipped=no-new-contracts")
		return
	}
	for _, contract := range newContracts {
		step("sim_revoke_"+contract, func() (string, error) {
			return client.RevokeEntrust("contractNo", asset, contract)
		})
	}
	finalOut, finalErr := step("sim_entrust_final", func() (string, error) {
		return client.GetEntrust(asset, "today", true)
	})
	needsFallback := finalErr != nil || containsAnyContract(finalOut, newContracts)
	if needsFallback && len(beforeContracts) == 0 {
		fmt.Println("poc_contract_revoke_incomplete=true")
		step("sim_revoke_all_fallback", func() (string, error) {
			return client.RevokeEntrust("allBuyAndSell", asset, "")
		})
		fallbackOut, fallbackErr := step("sim_entrust_final_after_fallback", func() (string, error) {
			return client.GetEntrust(asset, "today", true)
		})
		if fallbackErr != nil || containsAnyContract(fallbackOut, newContracts) {
			fmt.Println("poc_cleanup_failed=true")
			os.Exit(1)
		}
		return
	}
	if needsFallback {
		fmt.Println("poc_cleanup_incomplete=true")
		os.Exit(1)
	}
}

func step(label string, fn func() (string, error)) (string, error) {
	out, err := fn()
	if err == nil && isBusinessFailure(out) {
		err = fmt.Errorf("business failure")
	}
	fmt.Printf("%s_err=%v\n", label, err)
	fmt.Printf("%s_out=%s\n", label, out)
	return out, err
}

func printResult(label string, out string, err error) {
	fmt.Printf("%s_err=%v\n", label, err)
	fmt.Printf("%s_out=%s\n", label, out)
	if err != nil || isBusinessFailure(out) {
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

func extractContracts(out string) []string {
	matches := regexp.MustCompile(`\b\d{9,}\b`).FindAllString(out, -1)
	seen := map[string]bool{}
	var contracts []string
	for _, match := range matches {
		if seen[match] {
			continue
		}
		seen[match] = true
		contracts = append(contracts, match)
	}
	return contracts
}

func diffContracts(before []string, after []string) []string {
	seenBefore := map[string]bool{}
	for _, contract := range before {
		seenBefore[contract] = true
	}
	var result []string
	for _, contract := range after {
		if !seenBefore[contract] {
			result = append(result, contract)
		}
	}
	return result
}

func containsAnyContract(out string, contracts []string) bool {
	for _, contract := range contracts {
		if strings.Contains(out, contract) {
			return true
		}
	}
	return false
}

func isBusinessFailure(out string) bool {
	return strings.Contains(out, "failed")
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	os.Exit(2)
}
