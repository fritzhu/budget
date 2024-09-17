package transact

import "time"

type Transaction struct {
	Date        time.Time
	FromAccount string
	ToAccount   string
	Amount      float64
}
