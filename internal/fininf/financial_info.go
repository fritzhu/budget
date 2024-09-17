package fininf

import (
	"os"
	"time"

	"github.com/fritzhu/budget/internal/date"
	"github.com/fritzhu/budget/internal/transact"
	"gopkg.in/yaml.v3"
)

type FinancialInfo struct {
	Accounts []Account `yaml:"accounts"`
}

func LoadFinancialInfo(infile string) (*FinancialInfo, error) {
	var financialInfo FinancialInfo

	data, err := os.ReadFile(infile)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &financialInfo)
	if err != nil {
		return nil, err
	}

	return &financialInfo, nil
}

func (fin *FinancialInfo) GetPaydayAccount() *Account {
	for _, account := range fin.Accounts {
		if account.Payday != (Payday{}) {
			return &account
		}
	}
	panic("payday account not found")
}

func (fin *FinancialInfo) GetTransactions(from, to time.Time) []*transact.Transaction {
	from = date.Date(from)
	to = date.Date(to)

	act := fin.GetPaydayAccount()
	txs := act.Payday.GetTransactions(act, from, to, fin)

	for _, a := range fin.Accounts {
		for _, env := range a.Envelopes {
			etx := env.GetTransactions(a.Name, act, from, to)
			txs = append(txs, etx...)
		}
		for _, exp := range a.Expenses {
			etx := exp.GetTransactions(a.Name, act, from, to)
			txs = append(txs, etx...)
		}
	}

	return txs
}

func (fin *FinancialInfo) GetPaydayTransfers(d, from, to time.Time, includeLeftovers bool) map[string]float64 {
	xfers := map[string]float64{}

	act := fin.GetPaydayAccount()
	for _, a := range fin.Accounts {
		for _, env := range a.Envelopes {
			xfers[a.Name] += env.GetTopupAmount(d, act, from, to)
		}
		for _, exp := range a.Expenses {
			xfers[a.Name] += exp.GetTopupAmount(act.Payday.Frequency.Duration)
		}
	}

	if includeLeftovers {
		if act.Payday.TransferLeftoverTo != "" && act.Payday.TransferLeftoverTo != act.Name {
			t := float64(0)
			for _, v := range xfers {
				t += v
			}
			xfers[act.Payday.TransferLeftoverTo] += (act.Payday.Amount - t)
		}
	}

	return xfers
}

func (fin *FinancialInfo) GetTopupAmount(d, from, to time.Time) float64 {
	t := float64(0)
	xfers := fin.GetPaydayTransfers(d, from, to, false)
	for _, v := range xfers {
		t += v
	}
	return t
}
