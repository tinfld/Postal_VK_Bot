package main

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BOT_TOKEN     string
	IP_MAIL       string
	MAIL_USER     string
	PASSWORD_USER string
	CHAT_ID       string
}

func envLoad() (Config, error) {
	var conf Config

	_ = godotenv.Load()

	conf.BOT_TOKEN = os.Getenv("BOT_TOKEN")
	conf.IP_MAIL = os.Getenv("IP_MAIL")
	conf.MAIL_USER = os.Getenv("MAIL_USER")
	conf.PASSWORD_USER = os.Getenv("PASSWORD_USER")
	conf.CHAT_ID = os.Getenv("CHAT_ID")
	return conf, nil
}
