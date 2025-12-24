package main

import (
	"errors"
	"fmt"
)

type User struct {
	ID      string
	Name    string
	Balance float64
}

var ErrUserIsNotExist error = errors.New("user is not exist")
var ErrNotEnoughMoney error = errors.New("not enough money")

func (u *User) Deposit(amount float64) {
	u.Balance += amount
}
func (u *User) Withdraw(amount float64) error {
	if u.Balance < amount {
		return ErrNotEnoughMoney
	}
	u.Balance -= amount
	return nil
}

type Transaction struct {
	FromID string
	ToID   string
	Amount float64
}
type PaymentSystem struct {
	Users        map[string]*User
	Transactions []Transaction
}

func (p *PaymentSystem) AddUser(u *User) {
	p.Users[u.ID] = u
}
func (p *PaymentSystem) AddTransaction(t Transaction) {
	p.Transactions = append(p.Transactions, t)
}
func (p *PaymentSystem) ProcessingTransaction(t Transaction) error {
	if _, ok := p.Users[t.ToID]; !ok {
		return ErrUserIsNotExist
	}
	if _, ok := p.Users[t.FromID]; !ok {
		return ErrUserIsNotExist
	}
	err := p.Users[t.FromID].Withdraw(t.Amount)
	if err != nil {
		return err
	}
	p.Users[t.ToID].Deposit(t.Amount)
	return nil
}
func main() {
	paymentSystem := PaymentSystem{
		Users:        make(map[string]*User),
		Transactions: nil,
	}
	fmt.Println("Создаю UserID: 1 с балансом 1000")
	fmt.Println("Создаю UserID: 2 с балансом 500")
	user1 := User{
		ID:      "1",
		Name:    "Ivan",
		Balance: 1000,
	}
	user2 := User{
		ID:      "2",
		Name:    "Andrei",
		Balance: 500,
	}
	paymentSystem.AddUser(&user1)
	paymentSystem.AddUser(&user2)
	fmt.Println("Перевожу с UserID: 1 на UserID: 2 сумму в размере 200")
	fmt.Println("Перевожу с UserID: 2 на UserID: 1 сумму в размере 50")
	paymentSystem.AddTransaction(Transaction{
		FromID: "1",
		ToID:   "2",
		Amount: 200,
	})
	paymentSystem.AddTransaction(Transaction{
		FromID: "2",
		ToID:   "1",
		Amount: 50,
	})
	for _, trans := range paymentSystem.Transactions {
		err := paymentSystem.ProcessingTransaction(trans)
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("Итого")
	fmt.Println("У первого пользователя должно получиться 850, получилось", user1.Balance)
	fmt.Println("У второго пользователя должно получиться 650, получилось", user2.Balance)
}
