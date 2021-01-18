package bank

import (
	"fmt"
	"github.com/YoshikiShibata/gpl/ch09/ex01/bank"
	"testing"
)

func TestWithdrawNormal(t *testing.T) {
	bank.Deposit(200)
	ok := bank.Withdraw(100)
	if !ok {
		t.Error("Result is false, want true")
		return
	}
	ok = bank.Withdraw(100)
	if !ok {
		t.Error("Result is false, want true")
		return
	}
	ok = bank.Withdraw(100)
	if ok {
		t.Error("Result is true, want false")
		return
	}
	if bank.Balance() != 0 {
		t.Errorf("Result is %d, want 0", bank.Balance())
	}
	fmt.Println(bank.Balance())
}

func TestBank(t *testing.T) {
	done := make(chan struct{})

	go func() {
		bank.Deposit(200)
		fmt.Println("=", bank.Balance())
		done <- struct{}{}
	}()

	go func() {
		bank.Deposit(100)
		done <- struct{}{}
	}()

	<-done
	<-done

	if got, want := bank.Balance(), 300; got != want {
		t.Errorf("Balance = %d, want %d", got, want)
	}
}
