package main

import (
	"encoding/json"
	"fmt"
	"github.com/Spok95/bookgame/game"
	"html/template"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"path/filepath"
)

type Paragraph struct {
	Text    string   `json:"text"`
	Options []string `json:"options"`
}

type Player struct {
	Name     string
	Skill    string
	Dex      int
	Strength int
	Luck     int
	Honor    int
}

var story map[string]Paragraph
var player *game.Player

func main() {
	playerPath := "players/Kostya.json"
	if _, err := os.Stat(playerPath); err == nil {
		p, err := game.LoadPlayer("Kostya")
		if err == nil {
			player = p
			fmt.Println("✅ Игрок Kostya загружен автоматически")
		} else {
			fmt.Println("❌ Не удалось загрузить игрока Kostya:", err)
		}
	}
	// Загружаем story.json
	file, err := os.Open("data/story.json")
	if err != nil {
		log.Fatal("Ошибка при загрузке story.json:", err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&story)
	if err != nil {
		log.Fatal("Ошибка декодирования JSON:", err)
	}

	// Статика и обработчики
	http.Handle("/static", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", paragraphHandler)
	http.HandleFunc("/new", newGameHandler)
	http.HandleFunc("/save", savePlayerHandler)
	http.HandleFunc("/load", loadPlayerHandler)

	log.Println("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func paragraphHandler(w http.ResponseWriter, r *http.Request) {
	para := r.URL.Query().Get("para")
	if para == "" {
		para = "1"
	}

	p, ok := story[para]
	if !ok {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles(filepath.Join("templates", "paragraph.html"))
	if err != nil {
		http.Error(w, "Ошибка шаблона", 500)
		return
	}

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
		ImageURL:    "/static/images/" + para,
		MusicURL:    "/static/music/default.mp3",
		Player:      player,
		SaveSuccess: r.URL.Query().Get("save") == "ok",
	}

	if player == nil {
		http.Redirect(w, r, "/new", http.StatusSeeOther)
		return
	}
	player.CurrentPara = para
	tmpl.Execute(w, data)
}

func newGameHandler(w http.ResponseWriter, r *http.Request) {
	if player != nil {
		http.Redirect(w, r, "/?para="+player.CurrentPara, http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		skill := r.FormValue("skill")

		dex := randInt(1, 6) + 6
		str := randInt(1, 6) + randInt(1, 6) + 12
		luck := randInt(1, 6)

		player = game.NewPlayer(name, skill, dex, str, luck)

		http.Redirect(w, r, "/?para=1", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("templates/new_player.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func randInt(min, max int) int {
	return rand.IntN(max-min+1) + min
}

func savePlayerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := player.Save("player/" + player.Name + ".json")
	if err != nil {
		http.Error(w, "Ошибка сохранения игрока: "+err.Error(), http.StatusInternalServerError)
		return
	}

	para := player.CurrentPara
	http.Redirect(w, r, "/?para="+para+"&save=ok", http.StatusSeeOther)
}

func loadPlayerHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Не указано имя игрока", http.StatusBadRequest)
		return
	}

	loadedPlayer, err := game.LoadPlayer(name)
	if err != nil {
		http.Error(w, "Ошибка загрузки игрока: "+err.Error(), http.StatusInternalServerError)
		return
	}

	player = loadedPlayer

	// Переход на последний посещённый параграф
	http.Redirect(w, r, "/?para="+player.CurrentPara, http.StatusSeeOther)
}
