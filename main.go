package main

import (
	"os"
	"fmt"
	"time"
	"github.com/addyreal/heap-string"
	smtp "github.com/xhit/go-simple-mail/v2"
)

func usage() {
	fmt.Println("Usage:", os.Args[0], "<passwordfile> <host> <address> <from> <to> <subject> <message>")
}

func trim(b []byte) []byte {
	end := len(b)
	for end > 0 {
		last := b[end - 1]
		if last == '\n' || last == '\r' {
			end--
		} else {
			break
		}
	}

	return b[:end]
}

var Config config
type config struct {
	auth		string
	host		string
	from		string
	addr		string
}

func send(t string, s string, m string) error {
	server := smtp.NewSMTPClient()
	server.Host = Config.host
	server.Port = 587
	server.Username = Config.addr
	server.Password = Config.auth
	server.Encryption = smtp.EncryptionSTARTTLS
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	client, err := server.Connect()
	if err != nil {
		return err
	}

	email := smtp.NewMSG()
	email.SetFrom(fmt.Sprintf("%s <%s>", Config.from, Config.addr))
	email.AddTo(t)
	email.SetSubject(s)
	email.SetBody(smtp.TextPlain, m)
	if email.Error != nil {
		return email.Error
	}

	err = email.Send(client)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	args := os.Args[1:]
	if len(args) != 7 {
		usage()
		os.Exit(1)
	}

	err, buffer := heapstr.FromFile(args[0])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer buffer.Free()

	Config = config {
		auth:		string(trim(buffer.Get())),
		host:		args[1],
		addr:		args[2],
		from:		args[3],
	}

	err = send(args[4], args[5], args[6])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
