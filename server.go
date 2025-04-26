package main

import (
	"fmt"
	"github.com/Spok95/bookgame/game"
	"html/template"
	"log"
	rand "math/rand/v2"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var story *game.Story
var player *game.Player
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
	return player
}

func extractNextParas(text string) (victoryPara string, defeatPara string) {
	r := regexp.MustCompile(`/para\?para=(\d+)`)
	matches := r.FindAllStringSubmatch(text, -1)
	if len(matches) >= 2 {
		return matches[0][1], matches[1][1]
	} else if len(matches) == 1 {
		return matches[0][1], ""
	}
	return "", ""
}

func main() {
	// Загрузка сюжета
	var err error
	story, err = game.LoadStory("data/story.json")
	if err != nil {
		log.Fatal("Ошибка загрузки истории:", err)
	}
	log.Printf("✅ Загружено параграфов: %d\n", len(story.Paragraphs))

	// Попытка авто-загрузки игрока
	if p, err := game.LoadPlayer("Kostya"); err == nil {
		player = p
		fmt.Println("✅ Игрок Kostya загружен автоматически")
	}

	// Роуты
	http.HandleFunc("/", mainMenuHandler)
	http.HandleFunc("/new", newGameHandler)
	http.HandleFunc("/para", paragraphHandler)
	http.HandleFunc("/save", savePlayerHandler)
	http.HandleFunc("/load", loadPlayerHandler)
	http.HandleFunc("/players", listPlayersHandler)
	http.HandleFunc("/delete", deletePlayerHandler)
	http.HandleFunc("/fight", fightHandler)
	http.HandleFunc("/rules", rulesHandler)
	http.HandleFunc("/intro", introHandler)
	http.HandleFunc("/start", startHandler)

	// Статика
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Println("✅ Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func mainMenuHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/main_menu.html")
	tmpl.Execute(w, nil)
}

func introHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "intro.html", nil)
	if err != nil {
		http.Error(w, "Ошибка отображения предисловия: "+err.Error(), http.StatusInternalServerError)
	}
	if !hasItem(player.Inventory, "15 экю") {
		player.Inventory = append(player.Inventory, "15 экю")
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

func startHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/para?para=1", http.StatusSeeOther)
}

func newGameHandler(w http.ResponseWriter, r *http.Request) {
	player = nil
	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		skill := r.FormValue("skill")
		dex := randInt(1, 6) + 6
		str := randInt(1, 6) + randInt(1, 6) + 12
		luck := randInt(1, 6)

		player = game.NewPlayer(name, skill, dex, str, luck)
		http.Redirect(w, r, "/intro", http.StatusSeeOther)
		return
	}

	tmpl, _ := template.ParseFiles("templates/new_player.html")
	tmpl.Execute(w, nil)
}

func rulesHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "rules.html", nil)
	if err != nil {
		http.Error(w, "Ошибка отображения правил: "+err.Error(), http.StatusInternalServerError)
	}
}

func paragraphHandler(w http.ResponseWriter, r *http.Request) {
	p := getCurrentPlayer(r)
	if p == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	para := r.URL.Query().Get("para")
	if para == "" {
		para = p.CurrentPara
	}

	pg, ok := story.Paragraphs[para]
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

func savePlayerHandler(w http.ResponseWriter, r *http.Request) {
	if player == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	err := player.Save("players/" + player.Name + ".json")
	if err != nil {
		http.Error(w, "Ошибка сохранения игрока", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/para?para="+player.CurrentPara+"&save=ok", http.StatusSeeOther)
}

func loadPlayerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		p, err := game.LoadPlayer(name)
		if err != nil {
			renderLoadForm(w, "Игрок не найден: "+name)
			return
		}
		player = p
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
		player = p
		http.Redirect(w, r, "/para?para="+player.CurrentPara, http.StatusSeeOther)
		return
	} else {
		player = nil
	}

	renderLoadForm(w, "")
}

func deletePlayerHandler(w http.ResponseWriter, r *http.Request) {
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

func listPlayersHandler(w http.ResponseWriter, r *http.Request) {
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

func fightHandler(w http.ResponseWriter, r *http.Request) {
	p := getCurrentPlayer(r)
	if p == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	para := r.URL.Query().Get("para")
	if para == "" {
		para = p.CurrentPara
	}

	pg, ok := story.Paragraphs[para]
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
		victoryPara, defeatPara = extractNextParas(pg.Text)
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
