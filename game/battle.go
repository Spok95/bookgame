package game

import (
	"math/rand/v2"
	"regexp"
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

func RandInt(min, max int) int {
	return rand.IntN(max-min+1) + min
}

func ParseFightTag(text string) (Enemy, bool) {
	start := strings.Index(text, "fight:")
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

func ExtractNextParas(text string) (victoryPara string, defeatPara string) {
	r := regexp.MustCompile(`/para\?para=(\d+)`)
	matches := r.FindAllStringSubmatch(text, -1)
	if len(matches) >= 2 {
		return matches[0][1], matches[1][1]
	} else if len(matches) == 1 {
		return matches[0][1], ""
	}
	return "", ""
}
