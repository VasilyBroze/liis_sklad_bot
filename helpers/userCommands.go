package helpers

import (
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"liis_sklad_bot/entity"
	"strings"
	"time"
)

func ExecuteUserCommand(update tgbotapi.Update, database *sql.DB, BotSets entity.BotSettings, bot *tgbotapi.BotAPI, Stocks []entity.Stock) {
	ok := CheckRegisterUser(database, update)
	//ЕСЛИ ЗАБАНЕН - НИЧЕГО НЕ ДЕЛАЕМ
	if ok == false {
		return
	}
	command := strings.Split(update.Message.Text, " ")
	command[0] = strings.ToUpper(command[0])
	switch command[0] {

	case "ОСТАТОК": //ДОБАВИТЬ ЧЕЛОВЕКА В БД
		if len(command) != 1 {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Не понял Вас, Жду сообщения вида ОСТАТОК"))
		} else {
			var text string
			if len(Stocks) > 0 {
				fmt.Println(Stocks)
				counter := 0
				text = "На Складах сейчас доступны следующие позиции:\nНаименование, Артикул, Остаток, Единица\n"
				for _, position := range Stocks {

					//TODO: Мб убрать описание? закрепить его выше в начале сообщения
					text += fmt.Sprintf("✅<b>%v</b>, %v, %v, %v\n", position.StockName, position.Number, position.Quantity, position.Units)
					counter++
					//ОТПРАВЛЯЕМ СООБЩЕНИЕ КАЖДЫЕ 50 ПОЗИЦИЙ
					if counter%40 == 0 {
						fmt.Println("Попали в деление на 40")
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
						msg.ParseMode = "HTML"
						msg.ReplyMarkup = NumericKeyboard
						bot.Send(msg)
						time.Sleep(1 * time.Second)
						text = ""
					}
				}
				text += fmt.Sprintf("\nКоличество позиций: %v", counter)
			} else {
				text = "На складах нет позиций, либо возникла ошибка"
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
			msg.ParseMode = "HTML"
			msg.ReplyMarkup = NumericKeyboard
			bot.Send(msg)

		}

	case "/DESCRIPTION", "ОПИСАНИЕ":
		text := fmt.Sprintf("Описание команд:\nОстаток - выводит все доступные позиции" +
			"\nПоиск Кабель сетевой - поиск по артикулу/наименованию позиции\n")

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
		msg.ReplyMarkup = NumericKeyboard

		bot.Send(msg)

	case "/START", "START":
		//text := fmt.Sprintf("Привет! Я помогу тебе поздравлять твоих коллег без десятков надоедливых чатов :)\n" +
		//	"Для начала, зарегистрируйся. Напишите примерно так: \n*Регистрация Иван Иванов* (сначала имя, потом фамилия, Слово Регистрация тоже нужно писать)\n" +
		//	"Чтобы узнать что я умею введи Описание\n" +
		//	"Хорошего тебе дня!")
		text := fmt.Sprintf("Описание команд:\nОстаток - выводит все доступные позиции" +
			"\nПоиск Кабель сетевой - поиск по артикулу/наименованию позиции\n")
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
		msg.ReplyMarkup = NumericKeyboard
		bot.Send(msg)

	case "/SUBSCRIBERSTG", "ПОДПИСЧИКИТГ": //ВЫВОДИТ СПИСОК ВСЕХ ПОДПИСАВШИХСЯ
		text := ""
		var sum int
		//ЗАПРАШИВАЕМ ИЗ БД ВСЕ ИМЕНА
		rows, err := database.Query("SELECT id, username FROM users WHERE is_banned = false")
		if err != nil {
			fmt.Println("[ERROR] Error while got followers from DB: ", err)
		}

		var tgUsername string
		var chatID int64

		//ПРОВЕРЯЕМ ВСЕ ДАННЫЕ В БАЗЕ ИМЁН
		for rows.Next() {
			rows.Scan(&tgUsername, &chatID)

			sum += 1
			//ЕСЛИ ЮЗЕРНЕЙМ КОРРЕКТНЫЙ
			if tgUsername != "0" {
				if text != "" {
					text += fmt.Sprintf("\n@%s id - %v", tgUsername, chatID)
				} else {
					text += fmt.Sprintf("@%s id - %v", tgUsername, chatID)
				}
			} else {
				if text != "" {
					text += fmt.Sprintf("\n%v", chatID)
				} else {
					text += fmt.Sprintf("%v", chatID)
				}
			}

			if sum%100 == 0 {
				msg1 := tgbotapi.NewMessage(update.Message.Chat.ID, text)
				bot.Send(msg1)
				text = ""
			}
		}
		text += fmt.Sprintf("\nЧисло подписчиков: %v", sum)
		//ВЫВОД В ЧАТ
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
		msg.ReplyMarkup = NumericKeyboard
		bot.Send(msg)

	case "ПОИСК":

		if len(command) < 6 && len(command) > 1 { //ЕСЛИ ПОЛЬЗОВАТЕЛЬ ВВЕЛ КОМАНДУ ВЕРНО

			var inputName string
			_, inputName, _ = strings.Cut(update.Message.Text, " ")
			//Если введено хотябы 3 буквы
			if len(inputName) > 4 {
				//TODO: подумать, нужно ли заменять Е на Е
				//inputName = strings.Replace(inputName, "ё", "е", -1)
				//СЧЕТЧИК НАЙДЕННЫЙ СОВПАДЕНИЙ
				counter := 1
				text := "Вот что я нашел по запросу " + inputName + ":\n Наименование,Артикул, Остаток, Единица"
				textPos := ""
				for _, position := range Stocks {
					//peopleNameWithoutE := strings.Replace(people.Name, "ё", "е", -1)
					//TODO: Мб убрать описание? закрепить его выше в начале сообщения
					textPos = fmt.Sprintf("\n✅<b>%v</b>, %v, %v, %v", position.StockName, position.Number, position.Quantity, position.Units)
					if strings.Contains(strings.ToLower(textPos), strings.ToLower(inputName)) {
						counter++
						text += textPos
					}
					//ОТПРАВЛЯЕМ СООБЩЕНИЕ КАЖДЫЕ 50 ПОЗИЦИЙ
					if counter%50 == 0 {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
						msg.ParseMode = "HTML"
						msg.ReplyMarkup = NumericKeyboard
						bot.Send(msg)
						text = ""
					}
				}

				if strings.HasSuffix(text, "Остаток, Единица") == false {
					text += fmt.Sprintf("\nНайдено позиций: %v", counter-1)
					msg1 := tgbotapi.NewMessage(update.Message.Chat.ID, text)
					msg1.ParseMode = "HTML"
					msg1.ReplyMarkup = NumericKeyboard
					bot.Send(msg1)
				} else {
					msg1 := tgbotapi.NewMessage(update.Message.Chat.ID, "Совпадений не найдено")
					msg1.ReplyMarkup = NumericKeyboard
					bot.Send(msg1)
				}
			} else {
				msg1 := tgbotapi.NewMessage(update.Message.Chat.ID, "Чуть больше букв, пожалуйста :)")
				msg1.ReplyMarkup = NumericKeyboard
				bot.Send(msg1)
			}

		} else if len(command) == 1 { //ЕСЛИ НЕ ВВЕЛИ ИМЯ
			text := "Не понял, что искать, напишите примерно так - Поиск Кабель"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
			msg.ReplyMarkup = NumericKeyboard
			bot.Send(msg)
		} else {
			text := "Слишком много слов :( чуть меньше, пожалуйста"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
			msg.ReplyMarkup = NumericKeyboard
			bot.Send(msg)
		}

	case "ПАНИКА":
		if len(command) > 1 {
			if command[1] == "ПАНИКА" {
				//trim := fmt.Sprintf("%v %v ", command[0], BotSets.Anonce_pass)
				//msg := strings.TrimPrefix(update.Message.Text, trim)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Отключаюсь")
				bot.Send(msg)
				panic("Panic from Chat")
			}
		}

	default:
		text := "Команда не найдена, посмотрите Описание"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
		bot.Send(msg)
	}
}
