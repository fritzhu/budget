package main

import (
	"errors"
	"flag"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/fritzhu/budget/internal/date"
	"github.com/fritzhu/budget/internal/fininf"
	"github.com/fritzhu/budget/internal/transact"
)

var inputFile = flag.String("i", "./financial_info.yaml", "Input file containing your financial data, in YAML format.")

func main() {
	flag.Parse()

	fin, err := loadData()
	if err != nil {
		panic(err)
	}

	from := time.Now().AddDate(0, -6, 0)
	to := time.Now().AddDate(0, 6, 0)
	txs := fin.GetTransactions(from, to)

	ledger := transact.LedgerFromTransactions(txs)
	balances := ledger.CalculateMinimumBalancesAsToday(from, date.Date(time.Now()), to)

	fmt.Println("Minimum safe balances at EOD today:")
	for k, v := range balances {
		if k == "" {
			continue
		}
		fmt.Printf("\t%v = %v\n", k, math.Round(v*100.0)/100.0)
	}

	fmt.Println("\nTransfers to make this payday:")
	pda := fin.GetPaydayAccount()
	pd := pda.Payday
	is := date.NewIntervalStep(pd.LastPaid.Date, pd.Frequency.Duration, from, to)
	xfers := fin.GetPaydayTransfers(is.FirstOnOrAfter(date.Date(time.Now())), date.Date(from), date.Date(to), true)
	for k, v := range xfers {
		if k == "" || k == pda.Name {
			continue
		}
		fmt.Printf("\t => %v  %v\n", k, math.Round(v*100.0)/100.0)
	}
}

func loadData() (*fininf.FinancialInfo, error) {
	fi, err := os.Stat(*inputFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("input file does not exist")
		} else {
			return nil, err
		}
	}
	if fi.IsDir() {
		return nil, errors.New("input file is a directory")
	}

	fin, err := fininf.LoadFinancialInfo(*inputFile)
	if err != nil {
		return nil, err
	}

	return fin, err
}
