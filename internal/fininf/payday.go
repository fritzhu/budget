package fininf

import (
	"time"

	"github.com/fritzhu/budget/internal/date"
	"github.com/fritzhu/budget/internal/transact"
)

type Payday struct {
	Frequency          Interval       `yaml:"frequency"`
	Amount             float64        `yaml:"amount"`
	LastPaid           date.DateValue `yaml:"last_paid"`
	TransferLeftoverTo string         `yaml:"transfer_leftover_to"`
}

func (p *Payday) GetTransactions(paydayAccount *Account, from, to time.Time, fin *FinancialInfo) []*transact.Transaction {
	payday := date.NewIntervalStep(p.LastPaid.Date, p.Frequency.Duration, from, to)

	txs := []*transact.Transaction{}
	for d := payday.FirstOnOrAfter(date.Date(from)); d.Compare(to) <= 0; d = payday.FirstAfter(d) {
		if d == (time.Time{}) {
			break
		}
		tx := &transact.Transaction{
			Date:        d,
			FromAccount: "",
			ToAccount:   paydayAccount.Name,
			Amount:      p.Amount,
		}
		txs = append(txs, tx)

		amt := p.Amount - fin.GetTopupAmount(d, from, to)
		lot := p.TransferLeftoverTo
		if lot == "" {
			lot = paydayAccount.Name
		}
		if p.TransferLeftoverTo != paydayAccount.Name {
			tx = &transact.Transaction{
				Date:        d,
				FromAccount: paydayAccount.Name,
				ToAccount:   p.TransferLeftoverTo,
				Amount:      amt,
			}
			txs = append(txs, tx)
		}
		tx = &transact.Transaction{
			Date:        d,
			FromAccount: lot,
			ToAccount:   "",
			Amount:      amt,
		}
		txs = append(txs, tx)
	}

	return txs
}
