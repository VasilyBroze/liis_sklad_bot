package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"os"
	"time"

	//tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
	"liis_sklad_bot/entity"
	"liis_sklad_bot/helpers"
)

func main() {
	//СОЗДАНИЕ БД
	os.Create("./Stocks_Users.db")
	database, _ := sql.Open("sqlite3", "./Stocks_Users.db")
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, username TEXT, is_banned BOOLEAN)")
	if err != nil {
		fmt.Println(err)
	}
	statement.Exec()
	if err != nil {
		fmt.Println(err)
	}

	var BotSets entity.BotSettings

	//ИЗВЛЕКАЕМ ИЗ ФАЙЛА С НАСТРОЙКАМИ ПОЛЯ
	bs, err := helpers.GetSettings("settings.json")
	if err != nil {
		fmt.Println("open file error: " + err.Error())
		return
	}

	if err := json.Unmarshal(bs, &BotSets); err != nil {
		fmt.Println("setts")
		fmt.Println(err)

		return
	}

	bot, err := tgbotapi.NewBotAPI(BotSets.Bot_token)
	if err != nil {
		log.Panic(err)
	}
	Stocks := []entity.Stock{}
	// ОБНОВЛЕНИЕ АССОРТИМЕНТА
	go func() {
		for {
			resp, _ := http.Get(fmt.Sprintf("https://tools.aimylogic.com/api/googlesheet2json?sheet=%v&id=%v", BotSets.Google_sheet_stock_list, BotSets.Google_sheet_stock_url))

			Stocks = []entity.Stock{}

			err = json.NewDecoder(resp.Body).Decode(&Stocks)
			if err != nil {
				fmt.Println(resp.Body)
				fmt.Println(err, " body, err")
			}
			resp.Body.Close()

			for i, position := range Stocks {
				if position.Quantity == "" || position.Quantity == " " {
					Stocks[i].Quantity = "?"
				}
				if position.Number == "" || position.Number == " " {
					Stocks[i].Number = "?"
				}
				if position.Units == "" || position.Units == " " {
					Stocks[i].Units = "?"
				}
			}

			time.Sleep(2 * time.Hour) //hour
		}
	}()

	//НАСТРОЙКА СЛУШАТЕЛЯ
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	//ВЫПОЛНЕНИЕ КОМАНД ПОЛЬЗОВАТЕЛЯ
	for update := range updates {
		if update.Message != nil { // If we got a message
			helpers.ExecuteUserCommand(update, database, BotSets, bot, Stocks)
		} // else {
		//	if update.CallbackQuery != nil {
		//		helpers.ExecuteCallbackCommand(update, database, BotSets, bot, statement)
		//	}
		//}
	}
}
