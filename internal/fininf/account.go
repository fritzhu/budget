package fininf

type Account struct {
	Name      string     `yaml:"name"`
	Payday    Payday     `yaml:"payday"`
	Envelopes []Envelope `yaml:"envelopes"`
	Expenses  []Expense  `yaml:"expenses"`
}
