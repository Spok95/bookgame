package handlers

import (
	"github.com/Spok95/bookgame/game"
	"html/template"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var Story *game.Story
var Player *game.Player
var skipFightCheck bool
var remainingEnemies int
var victoryPara string
var defeatPara string
var templates = template.Must(template.New("").Funcs(template.FuncMap{
	"contains": contains,
}).ParseGlob("templates/*.html"))

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func getCurrentPlayer(r *http.Request) *game.Player {
	return Player
}

func MainMenuHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/main_menu.html")
	tmpl.Execute(w, nil)
}

func IntroHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "intro.html", nil)
	if err != nil {
		http.Error(w, "Ошибка отображения предисловия: "+err.Error(), http.StatusInternalServerError)
	}
	if !hasItem(Player.Inventory, "15 экю") {
		Player.Inventory = append(Player.Inventory, "15 экю")
	}
}

func hasItem(inv []string, item string) bool {
	for _, v := range inv {
		if v == item {
			return true
		}
	}
	return false
}

func StartHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/para?para=1", http.StatusSeeOther)
}

func NewGameHandler(w http.ResponseWriter, r *http.Request) {
	Player = nil
	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		skill := r.FormValue("skill")
		dex := randInt(1, 6) + 6
		str := randInt(1, 6) + randInt(1, 6) + 12
		luck := randInt(1, 6)

		Player = game.NewPlayer(name, skill, dex, str, luck)
		http.Redirect(w, r, "/intro", http.StatusSeeOther)
		return
	}

	tmpl, _ := template.ParseFiles("templates/new_player.html")
	tmpl.Execute(w, nil)
}

func RulesHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "rules.html", nil)
	if err != nil {
		http.Error(w, "Ошибка отображения правил: "+err.Error(), http.StatusInternalServerError)
	}
}

func ParagraphHandler(w http.ResponseWriter, r *http.Request) {
	p := getCurrentPlayer(r)
	if p == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	para := r.URL.Query().Get("para")
	if para == "" {
		para = p.CurrentPara
	}

	pg, ok := Story.Paragraphs[para]
	if !ok {
		http.Error(w, "Параграф не найден", http.StatusNotFound)
		return
	}

	p.CurrentPara = para

	if !skipFightCheck {
		for _, tag := range pg.Tags {
			if strings.HasPrefix(tag, "fight") {
				http.Redirect(w, r, "/fight?para="+para, http.StatusSeeOther)
				return
			}
		}
	}
	skipFightCheck = false

	data := struct {
		Player      *game.Player
		Number      string
		Text        template.HTML
		ImageURL    string
		MusicURL    string
		SaveSuccess bool
	}{
		Player:      p,
		Number:      para,
		Text:        template.HTML(pg.Text),
		ImageURL:    "/static/images/" + para + ".jpg",
		MusicURL:    "/static/music/" + para + ".mp3",
		SaveSuccess: r.URL.Query().Get("save") == "ok",
	}

	err := templates.ExecuteTemplate(w, "paragraph.html", data)
	if err != nil {
		http.Error(w, "Ошибка шаблона: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func SavePlayerHandler(w http.ResponseWriter, r *http.Request) {
	if Player == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	err := Player.Save("players/" + Player.Name + ".json")
	if err != nil {
		http.Error(w, "Ошибка сохранения игрока", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/para?para="+Player.CurrentPara+"&save=ok", http.StatusSeeOther)
}

func LoadPlayerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		p, err := game.LoadPlayer(name)
		if err != nil {
			renderLoadForm(w, "Игрок не найден: "+name)
			return
		}
		Player = p
		http.Redirect(w, r, "/para?para="+p.CurrentPara, http.StatusSeeOther)
		return
	}

	name := r.URL.Query().Get("name")
	if name != "" {
		p, err := game.LoadPlayer(name)
		if err != nil {
			http.Error(w, "Ошибка загрузки игрока: "+err.Error(), http.StatusInternalServerError)
			return
		}
		Player = p
		http.Redirect(w, r, "/para?para="+Player.CurrentPara, http.StatusSeeOther)
		return
	} else {
		Player = nil
	}

	renderLoadForm(w, "")
}

func DeletePlayerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешён", http.StatusSeeOther)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "Имя не указано", http.StatusBadRequest)
		return
	}

	err := os.Remove("players/" + name + ".json")
	if err != nil {
		http.Error(w, "Ошибка удаления игрока: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("🗑 Игрок %s удалён", name)
	http.Redirect(w, r, "/players", http.StatusSeeOther)
}

func renderLoadForm(w http.ResponseWriter, errMsg string) {
	tmpl, _ := template.ParseFiles("templates/load_player.html")
	tmpl.Execute(w, struct {
		Error string
	}{Error: errMsg})
}

func ListPlayersHandler(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir("players")
	if err != nil {
		http.Error(w, "Ошибка чтения players/", http.StatusInternalServerError)
		return
	}

	var names []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".json") {
			names = append(names, strings.TrimSuffix(f.Name(), ".json"))
		}
	}

	tmpl, _ := template.ParseFiles("templates/load_list.html")
	tmpl.Execute(w, struct {
		Names []string
	}{names})
}

func FightHandler(w http.ResponseWriter, r *http.Request) {
	p := getCurrentPlayer(r)
	if p == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	para := r.URL.Query().Get("para")
	if para == "" {
		para = p.CurrentPara
	}

	pg, ok := Story.Paragraphs[para]
	if !ok {
		http.Error(w, "Параграф не найден", http.StatusNotFound)
		return
	}

	if remainingEnemies == 0 {
		for _, tag := range pg.Tags {
			if strings.HasPrefix(tag, "fight") {
				parts := strings.Split(tag, ",")
				if len(parts) == 2 {
					num, err := strconv.Atoi(strings.TrimSpace(parts[1]))
					if err == nil {
						remainingEnemies = num
					}
				} else {
					remainingEnemies = 1
				}
			}
		}
		victoryPara, defeatPara = game.ExtractNextParas(pg.Text)
	}

	enemy := game.Enemy{
		Name: "Враг",
		Dex:  8,
		Str:  8,
	}
	// Вызов логики боя
	result := game.Fight(p, enemy)

	skipFightCheck = true

	data := map[string]interface{}{
		"Player":   p,
		"Enemy":    result.Enemy,
		"NextPara": victoryPara,
		"FailPara": defeatPara,
	}
	// Победа или поражение — показываем соответствующую страницу
	if result.Won {
		remainingEnemies--
		if remainingEnemies > 0 {
			http.Redirect(w, r, "/fight?para="+para, http.StatusSeeOther)
			return
		}
		remainingEnemies = 0
		tpl, _ := template.ParseFiles("templates/victory.html")
		tpl.Execute(w, data)
	} else {
		remainingEnemies = 0
		tpl, _ := template.ParseFiles("templates/defeat.html")
		tpl.Execute(w, data)
	}
}

func randInt(min, max int) int {
	return rand.IntN(max-min+1) + min
}
