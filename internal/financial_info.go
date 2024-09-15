package internal

type FinancialInfo struct {
	Payday   Payday     `yaml:"payday"`
	Envelope []Envelope `yaml:"envelopes"`
	Expense  []Expense  `yaml:"expenses"`
}
