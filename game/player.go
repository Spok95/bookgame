package game

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// Player — структура, описывающая персонажа игрока
type Player struct {
	Name      string
	Dexterity int    // Ловкость
	Strength  int    // Сила
	Luck      int    // Удача (1 = удачлив, 0 = неудачлив)
	Honor     int    // Честь
	Skill     string // Боевой навык
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
	fmt.Printf("Боевой навык: %s\n", p.Skill)
	fmt.Println("-------------------------------------")
}

// SelectSkill — предлагает выбрать боевой навык
func (p *Player) SelectSkill() {
	fmt.Println("\nВыберите боевое искусство (введите номер):")
	fmt.Println("1. Тайный удар шпагой")
	fmt.Println("2. Бой шпагой и кинжалом")
	fmt.Println("3. Стрельба из двух пистолетов")
	fmt.Println("4. Фехтование левой рукой")
	fmt.Println("5. Плавание")

	for {
		fmt.Print("Ваш выбор: ")
		var input string
		fmt.Scanln(&input)

		switch input {
		case "1":
			p.Skill = "Тайный удар шпагой"
		case "2":
			p.Skill = "Шпага и кинжал"
		case "3":
			p.Skill = "Два пистолета"
		case "4":
			p.Skill = "Фехтование левой рукой"
		case "5":
			p.Skill = "Плавание"
		default:
			fmt.Println("Некорректный выбор, попробуйте снова.")
			continue
		}
		break
	}
	fmt.Printf("Вы выбрали навык: %s\n", p.Skill)
}

// SaveToFile сохраняет игрока в JSON-файл
func (p *Player) SaveToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(p)
}

// LoadPlayer загружает игрока из JSON-файла
func LoadPlayer(filename string) (*Player, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var player Player
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&player)
	if err != nil {
		return nil, err
	}
	return &player, nil
}

// rollDice — бросок шестигранного кубика через новый генератор
func rollDice(r *rand.Rand) int {
	return r.Intn(6) + 1
}
