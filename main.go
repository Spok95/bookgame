package main

import (
	"bufio"
	"fmt"
	"github.com/Spok95/bookgame/game"
	"os"
	"strconv"
	"strings"
)

// Стартовая точка игры
func main() {
	fmt.Println("Добро пожаловать в 'Верную шпагу короля'")
	fmt.Println("___________________________________________")

	for {
		fmt.Println("\n1. Новая игра")
		fmt.Println("2. Загрузить игру")
		fmt.Println("3. Ввести номер параграфа")
		fmt.Println("0. Выход")
		fmt.Println("Выберите действие: ")

		input := readInput()

		switch input {
		case "1":
			fmt.Print("Введите имя персонажа: ")
			name := readInput()
			player := game.NewPlayer(name)
			player.Print()
		case "2":
			fmt.Println("Загрузка игры... (тоже заглушка)")
		case "3":
			fmt.Println("Введите номер параграфа: ")
			paraNumStr := readInput()
			paraNum, err := strconv.Atoi(paraNumStr)
			if err != nil {
				fmt.Println(err)
				continue
			}
			// Здесь пока фейковый текст
			fmt.Printf("\n[Параграф %d]:\n", paraNum)
			fmt.Println("Вы подъехали к воротам. Решите: войти (83) или вернуться (321)?")
		case "0":
			fmt.Println("Выход из игры. До встречи!")
			return
		default:
			fmt.Println("Неверный выбор. Повторите.")
		}
	}
}

// Функция чтения строки с консоли
func readInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
