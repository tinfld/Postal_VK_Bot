package main

import (
	botgolang "github.com/mail-ru-im/bot-golang"
)

func botVk(conf Config, logg Filename) (*botgolang.Bot, error) {
	// NewBot создаёт объект для взаимодействия с API бота
	bot, err := botgolang.NewBot(conf.BOT_TOKEN)
	if err != nil {
		logg.bot.Println("Ошибка подключения бота ", err)
		return nil, err
	}
	logg.bot.Println("Подключение к боту")

	return bot, nil
}
