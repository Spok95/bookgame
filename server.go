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

	// Статика
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Println("✅ Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func mainMenuHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/main_menu.html")
	tmpl.Execute(w, nil)
}

func newGameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		skill := r.FormValue("skill")
		dex := randInt(1, 6) + 6
		str := randInt(1, 6) + randInt(1, 6) + 12
		luck := randInt(1, 6)

		player = game.NewPlayer(name, skill, dex, str, luck)
		http.Redirect(w, r, "/para?para=1", http.StatusSeeOther)
		return
	}

	tmpl, _ := template.ParseFiles("templates/new_player.html")
	tmpl.Execute(w, nil)
}

func paragraphHandler(w http.ResponseWriter, r *http.Request) {
	if player == nil {
		http.Redirect(w, r, "/new", http.StatusSeeOther)
		return
	}

	para := r.URL.Query().Get("para")
	if para == "" {
		para = "1"
	}

	p, ok := story.Get(para)
	if !ok {
		http.NotFound(w, r)
		return
	}

	player.CurrentPara = para

	tmpl, _ := template.ParseFiles("templates/paragraph.html")
	data := struct {
		Number      string
		Text        string
		Option      []string
		ImageURL    string
		MusicURL    string
		Player      *game.Player
		SaveSuccess bool
	}{
		Number:      para,
		Text:        p.Text,
		Option:      p.Options,
		ImageURL:    "/static/images/" + para + ".jpg",
		MusicURL:    "/static/music/default.mp3",
		Player:      player,
		SaveSuccess: r.URL.Query().Get("save") == "ok",
	}
	tmpl.Execute(w, data)
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
	}

	renderLoadForm(w, "")
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

func randInt(min, max int) int {
	return rand.IntN(max-min+1) + min
}
