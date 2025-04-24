package game

import (
	"math/rand"
	"strconv"
	"strings"
)

type Enemy struct {
	Name string
	Dex  int
	Str  int
}

type FightRound struct {
	Round      int
	PlayerRoll int
	EnemyRoll  int
	Winner     string
	PlayerStr  int
	EnemyStr   int
}

type FightResult struct {
	Won       bool
	Log       []FightRound
	Enemy     Enemy
	ParaAfter string
}

func Fight(p *Player, enemy Enemy) FightResult {
	log := []FightRound{}
	round := 1

	for p.Strength > 0 && enemy.Str > 0 {
		playerRoll := RandInt(1, 6) + RandInt(1, 6) + p.Dex
		enemyRoll := RandInt(1, 6) + RandInt(1, 6) + enemy.Dex

		var winner string
		if playerRoll > enemyRoll {
			enemy.Str -= 2
			winner = p.Name
		} else if enemyRoll > playerRoll {
			p.Strength -= 2
			winner = enemy.Name
		} else {
			winner = "Ничья"
		}

		log = append(log, FightRound{
			Round:      round,
			PlayerRoll: playerRoll,
			EnemyRoll:  enemyRoll,
			Winner:     winner,
			PlayerStr:  p.Strength,
			EnemyStr:   enemy.Str,
		})
		round++
	}

	return FightResult{
		Won:   enemy.Str <= 0,
		Log:   log,
		Enemy: enemy,
	}
}

func ParseFightTag(text string) (Enemy, bool) {
	start := strings.Index(text, "#fight:")
	if start == -1 {
		return Enemy{}, false
	}

	tag := text[start+7:]
	parts := strings.Split(tag, ",")

	e := Enemy{}
	for _, part := range parts {
		pair := strings.Split(part, "=")
		if len(pair) != 2 {
			continue
		}
		key := strings.TrimSpace(pair[0])
		val := strings.TrimSpace(pair[1])

		switch key {
		case "Name":
			e.Name = val
		case "Dex":
			e.Dex = atoi(val)
		case "Str":
			e.Str = atoi(val)
		}
	}
	return e, true
}

func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func RandInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func (p *Player) AddItem(item string) bool {
	if len(p.Inventory) >= 5 {
		return false
	}
	for _, i := range p.Inventory {
		if i == item {
			return false
		}
	}
	p.Inventory = append(p.Inventory, item)
	return true
}
