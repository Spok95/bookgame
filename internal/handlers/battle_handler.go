package handlers

//type AttackResult struct {
//	PlayerRoll  int    `json:"PlayerRoll"`
//	EnemyRoll   int    `json:"EnemyRoll"`
//	Result      string `json:"Result"`
//	BattleEnded bool   `json:"BattleEnded"`
//}
//
//var EnemyHP = 10
//
//func AttackHandler(w http.ResponseWriter, r *http.Request) {
//	if r.Method != http.MethodPost {
//		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
//		return
//	}
//
//	playerRoll := rand.IntN(6) + 1
//	enemyRoll := rand.IntN(6) + 1
//
//	playerAttack := playerRoll + Player.Dex
//	enemyAttack := enemyRoll + 8
//
//	result := ""
//	battleEnded := false
//
//	if playerAttack > enemyAttack {
//		EnemyHP -= 2
//		result = "Вы нанесли урон противнику! (-2 HP)"
//	} else if playerAttack < enemyAttack {
//		Player.Strength -= 2
//		result = "Противник нанес вам урон! (-2 HP)"
//	} else {
//		result = "Парирование, без урона."
//	}
//
//	if EnemyHP <= 0 {
//		result += " 🎉 Победа!"
//		battleEnded = true
//	} else if Player.Strength <= 0 {
//		result += " ☠️ Вы проиграли!"
//		battleEnded = true
//	}
//
//	// (Пока без реального вычитания СИЛЫ — добавим позже.)
//	w.Header().Set("Content-Type", "application/json")
//	json.NewEncoder(w).Encode(AttackResult{
//		PlayerRoll:  playerRoll,
//		EnemyRoll:   enemyRoll,
//		Result:      result,
//		BattleEnded: battleEnded,
//	})
//}
