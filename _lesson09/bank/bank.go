package bank

type withdrawReq struct {
	amount  int
	resultc chan<- bool
}

var (
	deposits  = make(chan int)
	balances  = make(chan int)
	withdraws = make(chan *withdrawReq)
)

func Deposit(amount int) { deposits <- amount }
func Balance() int       { return <-balances }

func Withdraw(amount int) bool {
	resultc := make(chan bool)
	withdraws <- &withdrawReq{amount, resultc}
	return <-resultc
}

func teller() {
	var balance int
	for {
		select {
		case amount := <-deposits:
			balance += amount
		case balances <- balance:
		case req := <-withdraws:
			if balance >= req.amount {
				balance -= req.amount
				req.resultc <- true
			} else {
				req.resultc <- false
			}
		}
	}
}

func init() {
	go teller()
}
