package transact

import (
	"time"

	"github.com/fritzhu/budget/internal/date"
)

type Ledger struct {
	lines map[time.Time]*ledgerLine
}

type ledgerLine struct {
	deltas map[string]float64
	txs    []*Transaction
}

func LedgerFromTransactions(txs []*Transaction) *Ledger {
	l := Ledger{lines: map[time.Time]*ledgerLine{}}

	for _, tx := range txs {
		l.Apply(tx)
	}

	return &l
}

func (l *Ledger) CalculateMinimumBalancesAsToday(from, today, to time.Time) map[string]float64 {
	from, today, to = date.Date(from), date.Date(today), date.Date(to)

	balances := map[string]float64{}
	lowest := map[string]float64{}
	l.runTransactions(from, to, balances, lowest)

	balances = map[string]float64{}
	for k, v := range lowest {
		balances[k] = 0 - v
	}
	l.runTransactions(from, today, balances, lowest)

	return balances
}

func (i *Ledger) runTransactions(from, to time.Time, balances, lowest map[string]float64) {
	for d := date.Date(from); d.Compare(to) <= 0; d = date.Date(d.AddDate(0, 0, 1)) {
		if ll, ok := i.lines[d]; ok {
			for k, v := range ll.deltas {
				balances[k] += v
			}
			for k, v := range balances {
				if v < lowest[k] {
					lowest[k] = v
				}
			}
		}
	}
}

func (l *Ledger) Apply(tx *Transaction) {
	t := date.Date(tx.Date)
	if _, ok := l.lines[t]; !ok {
		l.lines[t] = &ledgerLine{deltas: map[string]float64{}, txs: []*Transaction{}}
	}

	l.lines[t].Apply(tx)
}

func (l *ledgerLine) Apply(tx *Transaction) {
	l.deltas[tx.FromAccount] -= tx.Amount
	l.deltas[tx.ToAccount] += tx.Amount
	l.txs = append(l.txs, tx)
}
