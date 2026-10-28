package main

import (
	"crypto/tls"
	"fmt"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
)

type mailConnect struct {
	ch      chan uint32 // Канал для уведомления об изменении числа писем в ящике
	options *imap.FetchOptions
	mailBox *imap.SelectData // Информация о выбранном почтовом ящике INBOX.
	client  *imapclient.Client
	section *imap.FetchItemBodySection
}

func mailPost(conf Config, logg Filename) (*mailConnect, error) {
	var err error

	mailCon := mailConnect{
		ch: make(chan uint32),
	}

	mailCon.client, err = imapclient.DialTLS(
		conf.IP_MAIL,
		&imapclient.Options{
			UnilateralDataHandler: &imapclient.UnilateralDataHandler{
				Mailbox: func(data *imapclient.UnilateralDataMailbox) {
					mailCon.ch <- *data.NumMessages
				},
			},
			TLSConfig: &tls.Config{
				// отключает проверку сертификата сервера
				InsecureSkipVerify: true,
			},
		},
	)
	if err != nil {
		logg.mail.Println("Ошибка подключения ", err)
		return &mailConnect{}, err
	}

	logg.mail.Println("Подключение к почте")
	// Авторизуемся в почтовом ящике
	err = mailCon.client.Login(
		conf.MAIL_USER,
		conf.PASSWORD_USER,
	).Wait()
	if err != nil {
		logg.mail.Println("Ошибка авторизации ", err)
		return &mailConnect{}, err
	}
	logg.mail.Println("Авторизация")
	// Выбираем папку inbox
	mailCon.mailBox, err = mailCon.client.Select("INBOX", nil).Wait()
	if err != nil {
		logg.mail.Println("Ошибка чтения ", err)
		return &mailConnect{}, err
	}
	if mailCon.mailBox.NumMessages == 0 {
		fmt.Println("Сообщений нет")
		return nil, nil
	}
	logg.mail.Println("Чтение почты")
	//получаем тело письма, UID и конверт с адресами и темой
	mailCon.section = &imap.FetchItemBodySection{}
	mailCon.options = &imap.FetchOptions{
		BodySection: []*imap.FetchItemBodySection{
			mailCon.section,
		},
		UID:      true,
		Envelope: true,
	}

	return &mailCon, nil
}
