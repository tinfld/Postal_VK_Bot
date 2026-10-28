package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-message/mail"
)

func main() {
	// Сначала открываем лог-файлы
	loger, err := CreateLog()
	if err != nil {
		fmt.Println("Ошибка создания логов ", err)
		return
	}
	// Закрываем открытые файлы логов при выходе из main
	defer loger.close()

	// Загружаем параметры из переменных окружения
	conf, err := envLoad()
	if err != nil {
		loger.app.Println("Ошибка чтения env ", err)
		return
	}

	// Создаём клиент VK-бота
	bot, err := botVk(conf, loger)
	if err != nil {
		loger.app.Println("Ошибка подключения бота ", err)
		return
	}

	// Подключаемся к IMAP и выбираем INBOX
	mailStruct, err := mailPost(conf, loger)
	if err != nil {
		loger.app.Println("Ошибка подключия к почте ", err)
		return
	}
	defer mailStruct.client.Close()
	defer func() {
		// Корректно завершаем IMAP-сеанс перед закрытием соединения
		err := mailStruct.client.Logout().Wait()
		if err != nil {
			loger.app.Println("Ошибка выхода", err)

		}
	}()

	fmt.Println("Старт работы")

	// Запоминаем число писем из выбранного ящика, при уведомлении IDLE оно будет обновлено
	numMsg := mailStruct.mailBox.NumMessages
	for {
		var txt string
		var txtBody string

		// IDLE позволяет ждать уведомлений сервера о новых письмах, не опрашивая его постоянно
		idleCmd, err := mailStruct.client.Idle()
		if err != nil {
			loger.app.Println("Ошибка ожидания ", err)
			return
		}

		// Ждём уведомление из обработчика Mailbox канал передаёт новое количество сообщений
		numMsg = <-mailStruct.ch

		// Прежде чем выполнить FETCH, завершаем команду IDLE
		if err := idleCmd.Close(); err != nil {
			loger.app.Println("Ошибка остановки ", err)
			return
		}

		// Получаем письмо по порядковому номеру и собираем результаты в срез
		messages, err := mailStruct.client.Fetch(imap.SeqSetNum(numMsg), mailStruct.options).Collect()
		if err != nil {
			loger.app.Println("Ошибка читаемости ", err)
			continue
		}

		for _, msg := range messages {
			body := msg.FindBodySection(mailStruct.section)
			if body == nil {
				continue
			}
			// Преобразуем байты тела письма
			reader := bytes.NewReader(body)

			mssg, err := mail.CreateReader(reader)
			if err != nil {
				loger.app.Println("Ошибка тела ", err)
				continue
			}

			for {
				part, err := mssg.NextPart()

				if err == io.EOF {
					break
				}

				if err != nil {
					loger.app.Println(err)
					break
				}
				// Проверяем Content-Type и сохраняем только обычный текст письма
				contentType := part.Header.Get("Content-Type")

				if strings.HasPrefix(contentType, "text/plain") {
					data, err := io.ReadAll(part.Body)
					if err != nil {
						loger.app.Println(err)
						continue
					}
					//fmt.Print(string(data))
					txtBody = string(data)
				}
			}

			from := msg.Envelope.From[0]
			to := msg.Envelope.To[0]

			//fmt.Println("\n", "От: ", from.Name, "\n", from.Mailbox+"@"+from.Host, "\n", "Кому: ", to.Name, "\n", to.Mailbox+"@"+to.Host, "\n", msg.Envelope.Date, "\n", msg.Envelope.Subject, "\n", "Тело: ", txtBody)
			txt += fmt.Sprintf("\nОт: %v \n %v@%v \n Кому: %v \n %v@%v \n %v \n %v \n Тело: %v", from.Name, from.Mailbox, from.Host, to.Name, to.Mailbox, to.Host, msg.Envelope.Date, msg.Envelope.Subject, txtBody)

			// тут можно сделать фильтр на домен
			if from.Host == "example" {
				botMsg := bot.NewTextMessage(conf.CHAT_ID, txt)
				err = botMsg.Send()
				if err != nil {
					loger.app.Println("Ошибка отправки ", err)
					continue
				}
				loger.app.Println("Отправка в чат")
			}
			// Это отдельная отправка, которая выполняется для каждого обработанного письма
			botMsg := bot.NewTextMessage("Адресс другого чата для одновременной отправки", txt)
			err = botMsg.Send()
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}

}
