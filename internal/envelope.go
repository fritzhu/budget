package internal

type Envelope struct {
	Name         string  `yaml:"name"`
	Account      string  `yaml:"account"`
	LastAmount   float64 `yaml:"last_amount"`
	LastDate     string  `yaml:"last_date"`
	TargetAmount float64 `yaml:"target_amount"`
	TargetDate   string  `yaml:"target_date"`
}
