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
func main() {
	user1 := User{
		ID:      "1",
		Name:    "Ivan",
		Balance: 570,
	}
	user2 := User{
		ID:      "2",
		Name:    "Andrei",
		Balance: 12360,
	}
	fmt.Println(user1.Name, "had", user1.Balance)
	user1.Deposit(5000)
	fmt.Println(user1.Name, "now has", user1.Balance)
	fmt.Println("Now we will try to withdraw more than one of users have")
	fmt.Println(user2.Name, "had", user2.Balance)
	err := user2.Withdraw(100000)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("we tried to withdraw too much, so we had an error")
}
