package game

import "fmt"

type Money struct {
	Coins int `json:"coins"`
	Sous  int `json:"sous"`
}

func (m *Money) ToSous() int {
	return m.Coins*30 + m.Sous
}

func (m *Money) Subtract(amount int) bool {
	total := m.ToSous()
	if total < amount {
		return false
	}
	total -= amount
	m.Coins = total / 30
	m.Sous = total % 30
	return true
}

func (m *Money) Add(amount int) {
	total := m.ToSous() + amount
	m.Coins = total / 30
	m.Sous = total % 30
}

func (m *Money) String() string {
	return fmt.Sprintf("%d экю %d су", m.Coins, m.Sous)
}
