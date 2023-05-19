package helpers

import (
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"os"
)

func GetSettings(path string) ([]byte, error) {

	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	a, err := io.ReadAll(f)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	f.Close()

	return a, err
}

func CheckRegisterUser(database *sql.DB, update tgbotapi.Update) bool {
	data1, _ := database.Query("SELECT id, username, is_banned FROM users WHERE id = ?", update.Message.Chat.ID)
	var isBanned bool
	var usernameDB string
	var idDB int64

	data1.Next()
	data1.Scan(&idDB, &usernameDB, &isBanned)
	data1.Close()

	//ЕСЛИ ЕЩЕ НЕ ПИСАЛ - РЕГИСТРИРУЕМ
	if idDB == 0 {
		statement, err := database.Prepare("INSERT INTO users (id, username, is_banned) VALUES (?, ?, ?)")
		if err != nil {
			fmt.Println(err)
		}
		statement.Exec(update.Message.Chat.ID, update.Message.From.UserName, false)
		if err != nil {
			fmt.Println(err)
		}
		return true
	} else {
		//ЕСЛИ ЗАБАНЕН - РЕТЕРНИМ 0 И ИГНОРИРУЕМ КОМАНДУ ПОЛЬЗОВАТЕЛЯ
		if isBanned == false {
			//ЕСЛИ НЕ ЗАБАНЕН ПРОВЕРЯЕМ ЮЗЕРНЕЙМ И ЕСЛИ НЕ СОВПАДАЕТ - ОБНОВЛЯЕМ
			//TODO: Дописать
			return true
		} else {
			return false
		}
	}
}
