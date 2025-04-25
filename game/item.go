package game

import "fmt"

type Item struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       Money  `json:"price"`
}

var StoreItems = []Item{
	{"Пистолет", "Оружие на 1 выстрел", Money{Coins: 4}},
	{"Шпага", "Для ближнего боя", Money{Coins: 3}},
	{"Пули и порох (5 выстрелов)", "Боеприпасы", Money{Coins: 3}},
	{"Кинжал", "На крайний случай", Money{Coins: 2}},
	{"Лошадь", "Может пригодиться", Money{Coins: 6}},
}

func (i Item) String() string {
	return fmt.Sprintf("%s — %s (%s)", i.Name, i.Description, i.Price.String())
}
