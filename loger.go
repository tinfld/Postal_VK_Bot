package main

import (
	"log"
	"os"
)

type Filename struct {
	app  *log.Logger
	mail *log.Logger
	bot  *log.Logger
	file []*os.File // Открытые файлы логов для последующего закрытия
}

func CreateLog() (Filename, error) {
	var fileLog Filename
	var fileName = [...]string{
		"logs/app.log",
		"logs/mail.log",
		"logs/bot.log",
	}

	for _, name := range fileName {
		file, err := os.OpenFile(
			name,
			os.O_CREATE|os.O_WRONLY|os.O_APPEND,
			0666,
		)
		if err != nil {
			return Filename{}, err
		}

		logger := log.New(file, "", log.Ldate|log.Ltime)

		if name == "logs/app.log" {
			fileLog.app = logger
		} else if name == "logs/mail.log" {
			fileLog.mail = logger
		} else if name == "logs/bot.log" {
			fileLog.bot = logger
		}
		fileLog.file = append(fileLog.file, file)
	}
	return fileLog, nil
}

// close закрывает все открытые файлы логов после завершения работы приложения
func (l Filename) close() {
	for _, file := range l.file {
		file.Close()
	}
}
