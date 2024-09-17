package fininf

import (
	"time"

	"github.com/fritzhu/budget/internal/date"
	"github.com/fritzhu/budget/internal/transact"
)

type Expense struct {
	Name      string         `yaml:"name"`                     // Name of the expense
	Amount    float64        `yaml:"amount"`                   // Amount of the expense
	Frequency Interval       `yaml:"frequency"`                // Payment frequency
	LastPaid  date.DateValue `yaml:"last_paid_date,omitempty"` // Last paid date
}

func (e *Expense) GetTransactions(accountName string, paydayAccount *Account, from, to time.Time) []*transact.Transaction {
	payday := date.NewIntervalStep(paydayAccount.Payday.LastPaid.Date, paydayAccount.Payday.Frequency.Duration, from, to)
	topupAmount := e.GetTopupAmount(paydayAccount.Payday.Frequency.Duration)

	txs := []*transact.Transaction{}
	d := e.Frequency.Duration
	if d == 0 {
		d = paydayAccount.Payday.Frequency.Duration
	}
	t := e.LastPaid.Date
	if t == (time.Time{}) {
		t = paydayAccount.Payday.LastPaid.Date
	}
	is := date.NewIntervalStep(t, d, from, to)
	for d := date.Date(from); d.Compare(to) <= 0; d = d.AddDate(0, 0, 1) {
		if is.IsOn(d) {
			tx := transact.Transaction{
				Date:        d,
				FromAccount: accountName,
				ToAccount:   "",
				Amount:      e.Amount,
			}
			if tx.FromAccount != tx.ToAccount {
				txs = append(txs, &tx)
			}
		}
		if payday.IsOn(d) {
			tx := transact.Transaction{
				Date:        d,
				FromAccount: paydayAccount.Name,
				ToAccount:   accountName,
				Amount:      topupAmount,
			}
			if tx.FromAccount != tx.ToAccount {
				txs = append(txs, &tx)
			}
		}
	}

	return txs
}

func (e *Expense) GetTopupAmount(payInterval time.Duration) float64 {
	d := e.Frequency.Duration
	if d == 0 {
		d = payInterval
	}
	return e.Amount / (float64(d) / float64(payInterval))
}
