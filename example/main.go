package main

import (
	"io/ioutil"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/scorredoira/email"
	"github.com/zhnxin/emailagent"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	contentType = kingpin.Arg("type", "content type").Required().Enum("plain", "html")
	confentFile = kingpin.Arg("content-file", "content body").Required().ExistingFile()
	configfile  = kingpin.Flag("conf", "authuration configuraion").Required().Short('c').ExistingFile()
	attachments = kingpin.Flag("attach", "attach file").Short('a').ExistingFiles()
	to          = kingpin.Flag("to", "target").Required().Short('t').Strings()
	subject     = kingpin.Flag("subject", "email subject").Required().Short('s').String()
	Cc          = kingpin.Flag("cc", "email cc").Strings()
	Bcc         = kingpin.Flag("bcc", "email bcc").Strings()
)

type EmailConfig struct {
	Host     string
	Port     int
	IsSSL    bool
	User     string
	NickName string
	Password string
}

func main() {
	kingpin.Parse()
	config := EmailConfig{}
	_, err := toml.DecodeFile(*configfile, &config)
	if err != nil {
		log.Fatal(err)
	}
	var msg *email.Message
	body, err := ioutil.ReadFile(*confentFile)
	if err != nil {
		log.Fatal(err)
	}
	switch *contentType {
	case "plain":
		msg = email.NewMessage(*subject, string(body))
	case "html":
		msg = email.NewHTMLMessage(*subject, string(body))
	default:
		log.Fatal("only support content type -- plain and html")
	}
	for _, a := range *attachments {
		if err := msg.Attach(a); err != nil {
			log.Fatal(err)
		}
	}
	msg.From.Address = config.User
	if config.NickName != "" {
		msg.From.Name = config.NickName
	} else {
		msg.From.Name = config.User
	}
	msg.To = *to
	msg.Cc = *Cc
	msg.Bcc = *Bcc
	agent := emailagent.New(config.User, config.Password, config.Host, config.Port, config.IsSSL)
	if err := agent.SendEmail(msg); err != nil {
		log.Fatal(err)
	}
	log.Println("success")

}
