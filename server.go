package main

import (
	"fmt"
	"github.com/Spok95/bookgame/game"
	"github.com/Spok95/bookgame/internal/handlers"
	"log"
	"net/http"
)

func main() {
	// Загрузка сюжета
	var err error
	handlers.Story, err = game.LoadStory("data/story.json")
	if err != nil {
		log.Fatal("Ошибка загрузки истории:", err)
	}
	log.Printf("✅ Загружено параграфов: %d\n", len(handlers.Story.Paragraphs))

	// Попытка авто-загрузки игрока
	if p, err := game.LoadPlayer("Kostya"); err == nil {
		handlers.Player = p
		fmt.Println("✅ Игрок Kostya загружен автоматически")
	}

	// Роуты
	http.HandleFunc("/", handlers.MainMenuHandler)
	http.HandleFunc("/new", handlers.NewGameHandler)
	http.HandleFunc("/para", handlers.ParagraphHandler)
	http.HandleFunc("/save", handlers.SavePlayerHandler)
	http.HandleFunc("/load-list", handlers.ListPlayersHandler)
	http.HandleFunc("/load-player", handlers.LoadFromListHandler)
	http.HandleFunc("/delete", handlers.DeletePlayerHandler)
	http.HandleFunc("/fight", handlers.FightHandler)
	http.HandleFunc("/rules", handlers.RulesHandler)
	http.HandleFunc("/intro", handlers.IntroHandler)
	http.HandleFunc("/start", handlers.StartHandler)
	http.HandleFunc("/roll", handlers.RollDiceHandler)
	http.HandleFunc("/attack", handlers.AttackHandler)
	http.HandleFunc("/luck", handlers.LuckHandler)

	// Статика
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Println("✅ Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
