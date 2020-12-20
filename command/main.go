package main

import (
	"fmt"
)

type BankAccount_ struct {
	Balance int
}

func Deposit_(ba *BankAccount_, amount int) {
	fmt.Println("Depositing", amount)
	ba.Balance += amount
}

func Withdraw_(ba *BankAccount_, amount int) {
	if ba.Balance >= amount {
		fmt.Println("Withdrawing", amount)
		ba.Balance -= amount
	}
}

var overdraftLimit = -500

type BankAccount struct {
	balance int
}

func (b *BankAccount) Deposit(amount int) {
	b.balance += amount
	fmt.Println("Deposited", amount, "\b, balance is now", b.balance)
}

func (b *BankAccount) WithDraw(amount int) bool {
	if b.balance-amount >= overdraftLimit {
		b.balance -= amount
		fmt.Println("Withdraw", amount, "\b, balance is now", b.balance)
		return true
	}
	return false
}

type Command interface {
	Call()
	Undo()
	Succeeded() bool
	SetSucceeded(value bool)
}

type Action int

const (
	Deposit Action = iota
	Withdraw
)

type BankAccountCommand struct {
	account   *BankAccount
	action    Action
	amount    int
	succeeded bool
}

func (b *BankAccountCommand) Succeeded() bool {
	return b.succeeded
}

func (b *BankAccountCommand) SetSucceeded(value bool) {
	b.succeeded = value
}

func (b *BankAccountCommand) Call() {
	switch b.action {
	case Deposit:
		b.account.Deposit(b.amount)
		b.succeeded = true
	case Withdraw:
		b.succeeded = b.account.WithDraw(b.amount)
	}
}

func (b *BankAccountCommand) Undo() {
	if !b.succeeded {
		return
	}
	switch b.action {
	case Deposit:
		b.account.WithDraw(b.amount)
	case Withdraw:
		b.account.Deposit(b.amount)
	}
}

func NewBankAccountCommand(account *BankAccount, action Action, amount int) *BankAccountCommand {
	return &BankAccountCommand{
		account: account,
		action:  action,
		amount:  amount,
	}
}

type CompositeBankAccountCommand struct {
	commands []Command
}

func (c *CompositeBankAccountCommand) Succeeded() bool {
	for _, cmd := range c.commands {
		if !cmd.Succeeded() {
			return false
		}
	}
	return true
}

func (c *CompositeBankAccountCommand) SetSucceeded(value bool) {
	for _, cmd := range c.commands {
		cmd.SetSucceeded(value)
	}
}

func (c *CompositeBankAccountCommand) Call() {
	for _, cmd := range c.commands {
		cmd.Call()
	}
}

func (c *CompositeBankAccountCommand) Undo() {
	for idx := range c.commands {
		c.commands[len(c.commands)-idx-1].Undo()
	}
}

type MoneyTransferCommand struct {
	CompositeBankAccountCommand
	from, to *BankAccount
	amount   int
}

func NewMoneyTransferCommand(from *BankAccount, to *BankAccount, amount int) *MoneyTransferCommand {
	c := &MoneyTransferCommand{
		from:   from,
		to:     to,
		amount: amount,
	}
	c.commands = append(c.commands, NewBankAccountCommand(from, Withdraw, amount))
	c.commands = append(c.commands, NewBankAccountCommand(to, Deposit, amount))
	return c
}

func (m *MoneyTransferCommand) Call() {
	ok := true
	for _, cmd := range m.commands {
		if ok {
			cmd.Call()
			ok = cmd.Succeeded()
		} else {
			cmd.SetSucceeded(false)
		}
	}
}

func main() {
	ba := &BankAccount{}
	cmdDeposit := NewBankAccountCommand(ba, Deposit, 100)
	cmdWithDraw := NewBankAccountCommand(ba, Withdraw, 1000)
	cmdDeposit.Call()
	cmdWithDraw.Call()
	fmt.Println(ba)
	cmdWithDraw.Undo()
	cmdDeposit.Undo()
	fmt.Println(ba)

	from := BankAccount{100}
	to := BankAccount{0}
	mtc := NewMoneyTransferCommand(&from, &to, 100)
	mtc.Call()

	fmt.Println("from=", from, "to=", to)

	fmt.Println("Undoing...")
	mtc.Undo()
	fmt.Println("from=", from, "to=", to)
}
