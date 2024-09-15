package internal

import "time"

type Expense struct {
	Name      string    `yaml:"name"`                     // Name of the expense
	Amount    float64   `yaml:"amount"`                   // Amount of the expense
	Account   string    `yaml:"account"`                  // Account where the payment is made
	Frequency string    `yaml:"frequency"`                // Payment frequency
	LastPaid  time.Time `yaml:"last_paid_date,omitempty"` // Last paid date
}
