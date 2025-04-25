package main

import (
	"fmt"
	"github.com/Spok95/bookgame/game"
	"html/template"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"strings"
)

var story *game.Story
var player *game.Player
var templates = template.Must(template.New("").Funcs(template.FuncMap{
	"contains": contains,
}).ParseGlob("templates/*.html"))

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
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
	if player == nil {
		http.Redirect(w, r, "/menu", http.StatusSeeOther)
		return
	}

	para := r.URL.Query().Get("para")
	if para == "" {
		para = player.CurrentPara
	}

	p, ok := story.Paragraphs[para]
	if !ok {
		http.Error(w, "Параграф не найден", http.StatusNotFound)
		return
	}

	player.CurrentPara = para

	data := struct {
		Player      *game.Player
		Number      string
		Text        template.HTML
		ImageURL    string
		MusicURL    string
		SaveSuccess bool
	}{
		Player:      player,
		Number:      para,
		Text:        template.HTML(p.Text),
		ImageURL:    "/static/images/" + para + ".jpg",
		MusicURL:    "/static/music/" + para + ".mp3",
		SaveSuccess: r.URL.Query().Get("save") == "ok",
	}

	if strings.Contains(p.Text, "#fight:") {
		http.Redirect(w, r, "/fight?para="+para, http.StatusSeeOther)
		return
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
	if player == nil {
		http.Redirect(w, r, "/new", http.StatusSeeOther)
		return
	}

	para := r.URL.Query().Get("para")
	if para == "" {
		para = player.CurrentPara
	}

	p, ok := story.Paragraphs[para]
	if !ok {
		http.NotFound(w, r)
		return
	}

	enemy, found := game.ParseFightTag(p.Text)
	if !found {
		http.Error(w, "В параграфе нет информации о враге", http.StatusNotFound)
		return
	}

	result := game.Fight(player, enemy)
	player.CurrentPara = result.ParaAfter

	data := struct {
		Player *game.Player
		game.FightResult
	}{
		Player:      player,
		FightResult: result,
	}

	err := templates.ExecuteTemplate(w, "fight.html", data)
	if err != nil {
		http.Error(w, "Ошибка шаблона: "+err.Error(), http.StatusInternalServerError)
	}
}

func randInt(min, max int) int {
	return rand.IntN(max-min+1) + min
}
