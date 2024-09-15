package internal

import "time"

type Payday struct {
	Frequency string    `yaml:"frequency"`
	Amount    float64   `yaml:"amount"`
	Account   string    `yaml:"account"`
	LastPaid  time.Time `yaml:"last_paid"`
}
