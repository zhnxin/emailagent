package emailagent

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"

	"github.com/scorredoira/email"
)

//Agent simple email client support ssl
type Agent struct {
	user  string
	addr  string
	isSSL bool
	auth  smtp.Auth
}

//New new email client
func New(user, password, host string, port int, isSsl bool) *Agent {
	ec := &Agent{
		user:  user,
		addr:  fmt.Sprintf("%s:%d", host, port),
		isSSL: isSsl,
		auth:  smtp.PlainAuth("", user, password, host),
	}
	return ec
}

func NewWithIdentify(identity, user, password, host string, port int, isSsl bool) *Agent {
	ec := &Agent{
		user:  user,
		addr:  fmt.Sprintf("%s:%d", host, port),
		isSSL: isSsl,
		auth:  smtp.PlainAuth(identity, user, password, host),
	}
	return ec
}

func (ec *Agent) sendMailTLS(msg *email.Message) error {
	host, _, _ := net.SplitHostPort(ec.addr)
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}
	conn, err := tls.Dial("tcp", ec.addr, tlsconfig)
	if err != nil {
		return fmt.Errorf("DialConn:%v", err)
	}
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("Agent:generateAgent:%v", err)
	}
	defer client.Close()
	if ec.auth != nil {
		if ok, _ := client.Extension("AUTH"); ok {
			if err = client.Auth(ec.auth); err != nil {
				return fmt.Errorf("Agent:clientAuth:%v", err)
			}
		}
	}
	if err = client.Mail(ec.user); err != nil {
		return fmt.Errorf("Agent:clientMail:%v", err)
	}

	for _, addr := range msg.Tolist() {
		if err = client.Rcpt(addr); err != nil {
			return fmt.Errorf("Agent:Rcpt:%v", err)
		}
	}
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("Agent:%v", err)
	}
	_, err = w.Write(msg.Bytes())
	if err != nil {
		return fmt.Errorf("Agent:WriterBody:%v", err)
	}
	err = w.Close()
	if err != nil {
		return fmt.Errorf("Agent:CloseBody:%v", err)
	}
	return client.Quit()
}

func (ec *Agent) sendMail(msg *email.Message) error {
	return smtp.SendMail(ec.addr, ec.auth, ec.user, msg.Tolist(), msg.Bytes())
}

//SendEmail send email by string content
func (ec *Agent) SendEmail(msg *email.Message) error {
	if ec.isSSL {
		return ec.sendMailTLS(msg)
	}
	return ec.sendMail(msg)
}
