package game

import (
	"fmt"
	"math/rand"
	"time"
)

// Player — структура, описывающая персонажа игрока
type Player struct {
	Name      string
	Dexterity int // Ловкость
	Strength  int // Сила
	Luck      int // Удача (1 = удачлив, 0 = неудачлив)
	Honor     int // Честь
}

// NewPlayer — создаёт нового игрока, бросая кубики по таблице
func NewPlayer(name string) *Player {
	// Создаём локальный генератор случайных чисел
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Таблицы соответствия бросков
	dexTable := map[int]int{
		1: 12,
		2: 8,
		3: 10,
		4: 7,
		5: 9,
		6: 11,
	}

	strTable := map[int]int{
		1: 22,
		2: 18,
		3: 14,
		4: 24,
		5: 16,
		6: 20,
	}

	dexRoll := rollDice(r)
	strRoll := rollDice(r)

	return &Player{
		Name:      name,
		Dexterity: dexTable[dexRoll],
		Strength:  strTable[strRoll],
		Luck:      1,
		Honor:     3,
	}
}

// Print — вывод параметров игрока
func (p *Player) Print() {
	fmt.Println("------ Листок путешественника ------")
	fmt.Printf("Имя: %s\n", p.Name)
	fmt.Printf("Ловкость: %d\n", p.Dexterity)
	fmt.Printf("Сила: %d\n", p.Strength)
	fmt.Printf("Удача: %d\n", p.Luck)
	fmt.Printf("Честь: %d\n", p.Honor)
	fmt.Println("-------------------------------------")
}

// rollDice — бросок шестигранного кубика через новый генератор
func rollDice(r *rand.Rand) int {
	return r.Intn(6) + 1
}
