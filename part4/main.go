package main

import (
	"errors"
	"fmt"
	"sync"
)

type User struct {
	ID      string
	Name    string
	Balance float64
	mu      sync.Mutex
}

var ErrUserIsNotExist error = errors.New("user is not exist")
var ErrNotEnoughMoney error = errors.New("not enough money")

func (u *User) Deposit(amount float64) {
	u.mu.Lock()
	u.Balance += amount
	u.mu.Unlock()
}
func (u *User) Withdraw(amount float64) error {
	u.mu.Lock()
	if u.Balance < amount {
		u.mu.Unlock()
		return ErrNotEnoughMoney
	}
	u.Balance -= amount
	u.mu.Unlock()
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
	wg := &sync.WaitGroup{}
	paymentSystem := PaymentSystem{
		Users:        make(map[string]*User),
		Transactions: nil,
	}
	fmt.Println("Создаю UserID: 1 с балансом 1000")
	fmt.Println("Создаю UserID: 2 с балансом 500")
	user1 := User{
		ID:      "1",
		Name:    "Ivan",
		Balance: 10000,
	}
	user2 := User{
		ID:      "2",
		Name:    "Andrei",
		Balance: 10000,
	}
	paymentSystem.AddUser(&user1)
	paymentSystem.AddUser(&user2)
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
	paymentSystem.AddTransaction(Transaction{
		FromID: "1",
		ToID:   "2",
		Amount: 100,
	})
	paymentSystem.AddTransaction(Transaction{
		FromID: "1",
		ToID:   "2",
		Amount: 200,
	})
	paymentSystem.AddTransaction(Transaction{
		FromID: "2",
		ToID:   "1",
		Amount: 500,
	})
	paymentSystem.AddTransaction(Transaction{
		FromID: "1",
		ToID:   "2",
		Amount: 1000,
	})
	paymentSystem.AddTransaction(Transaction{
		FromID: "2",
		ToID:   "1",
		Amount: 50,
	})
	paymentSystem.AddTransaction(Transaction{
		FromID: "1",
		ToID:   "2",
		Amount: 3000,
	})
	paymentSystem.AddTransaction(Transaction{
		FromID: "2",
		ToID:   "1",
		Amount: 150,
	})
	ch := make(chan Transaction, len(paymentSystem.Transactions))
	for _, tr := range paymentSystem.Transactions {
		ch <- tr
	}
	close(ch)
	for i := range 5 {
		wg.Add(1)
		go func() {
			fmt.Println("worker", i, "in process")
			err := paymentSystem.worker(wg, ch)
			if err != nil {
				fmt.Println(err)
			}
		}()
	}
	wg.Wait()
	fmt.Println(user1.Balance, user2.Balance)
}
func (p *PaymentSystem) worker(wg *sync.WaitGroup, ch <-chan Transaction) error {
	defer wg.Done()
	for t := range ch {
		err := p.ProcessingTransaction(t)
		if err != nil {
			return err
		}
	}
	return nil
}
