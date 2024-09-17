package fininf

import (
	"math"
	"time"

	"github.com/fritzhu/budget/internal/date"
	"github.com/fritzhu/budget/internal/transact"
)

type Envelope struct {
	Name         string         `yaml:"name"`
	StartAmount  float64        `yaml:"start_amount"`
	StartDate    date.DateValue `yaml:"start_date"`
	TargetAmount float64        `yaml:"target_amount"`
	TargetDate   date.DateValue `yaml:"target_date"`
}

func (e *Envelope) GetTransactions(accountName string, paydayAccount *Account, from, to time.Time) []*transact.Transaction {
	payday := date.NewIntervalStep(paydayAccount.Payday.LastPaid.Date, paydayAccount.Payday.Frequency.Duration, from, to)
	paydays := e.getPaydaysBetween(payday, e.StartDate.Date, e.TargetDate.Date)
	topupAmount := math.Round(((e.TargetAmount-e.StartAmount)/float64(paydays))*100.0) / 100.0
	/*fmt.Printf("\tEnv %v paydays %v\n", e.Name, paydays)
	fmt.Printf("\tEnv %v topup %v\n", e.Name, topupAmount)
	fmt.Printf("\tEnv %v target %v to %v\n", e.Name, e.TargetDate.Date, to)*/

	txs := []*transact.Transaction{}

	for d := payday.FirstOnOrAfter(from); d.Compare(to) <= 0; d = payday.FirstAfter(d) {
		if d == (time.Time{}) {
			break
		}
		if d.Compare(e.StartDate.Date) >= 0 && d.Compare(e.TargetDate.Date) <= 0 {
			tx := transact.Transaction{
				Date:        d,
				FromAccount: paydayAccount.Name,
				ToAccount:   accountName,
				Amount:      topupAmount,
			}
			txs = append(txs, &tx)
		}
	}

	tx := transact.Transaction{
		Date:        e.TargetDate.Date,
		FromAccount: accountName,
		ToAccount:   "",
		Amount:      e.TargetAmount,
	}
	txs = append(txs, &tx)

	return txs
}

func (e *Envelope) GetTopupAmount(d time.Time, paydayAccount *Account, from, to time.Time) float64 {
	d, from, to = date.Date(d), date.Date(from), date.Date(to)
	if d.Compare(e.StartDate.Date) >= 0 && d.Compare(e.TargetDate.Date) <= 0 {
		payday := date.NewIntervalStep(paydayAccount.Payday.LastPaid.Date, paydayAccount.Payday.Frequency.Duration, from, to)
		paydays := e.getPaydaysBetween(payday, e.StartDate.Date, e.TargetDate.Date)
		return math.Round(((e.TargetAmount-e.StartAmount)/float64(paydays))*100.0) / 100.0
	}
	return 0
}

func (e *Envelope) getPaydaysBetween(payday *date.IntervalStep, from, to time.Time) int {
	count := 0
	d := date.Date(from)
	for d.Compare(to) <= 0 {
		if payday.IsOn(d) {
			count++
		}
		d = d.AddDate(0, 0, 1)
	}
	return count
}
