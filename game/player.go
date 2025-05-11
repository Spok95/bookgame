package game

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Player — структура, описывающая персонажа игрока
type Player struct {
	Name        string   `json:"name"`
	Skill       string   `json:"skill"`
	Dex         int      `json:"dex"`
	Strength    int      `json:"strength"`
	Luck        int      `json:"luck"`
	Honor       int      `json:"honor"`
	CurrentPara string   `json:"current_para"`
	Inventory   []string `json:"inventory"`
	Money       Money    `json:"money"`
	LuckStreak  int      `json:"luck_streak"`
}

// NewPlayer — создаёт нового игрока
func NewPlayer(name, skill string, dex, strength, luck int) *Player {
	return &Player{
		Name:        name,
		Skill:       skill,
		Dex:         dex,
		Strength:    strength,
		Luck:        luck,
		Honor:       3,
		CurrentPara: "1",
		Inventory:   make([]string, 0),
		Money:       Money{Coins: 15, Sous: 0},
		LuckStreak:  0,
	}
}

// SaveToFile сохраняет игрока в JSON-файл
func (p *Player) Save(path string) error {
	_ = os.MkdirAll("players", 0755)
	filePath := filepath.Join("players", p.Name+".json")
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

// LoadPlayer загружает игрока по имени
func LoadPlayer(name string) (*Player, error) {
	filePath := filepath.Join("players", name+".json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var p Player
	err = json.Unmarshal(data, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// DebugPrint выводит игрока в консоль (опционально)
func (p *Player) DebugPrint() {
	fmt.Println("------ Листок путешественника ------")
	fmt.Println("Имя:", p.Name)
	fmt.Println("Навык:", p.Skill)
	fmt.Println("Ловкость:", p.Dex)
	fmt.Println("Сила:", p.Strength)
	fmt.Println("Удача:", p.Luck)
	fmt.Println("Честь:", p.Honor)
	fmt.Println("Экю:", p.Money.Coins, "| Су:", p.Money.Sous)
	fmt.Println("Параграф:", p.CurrentPara)
	fmt.Println("------------------------------------")
}
