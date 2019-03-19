# email

An easy way to send emails By SSL with `github.com/scorredoira/email`.

# install

```
go get github.com/zhnxin/emailagent
```

# Usage

You can refer to [github.com/scorredoira/email](https://github.com/scorredoira/email) to get more infomations about email message generation.

```go
package emailagent_test

import (
	"net/mail"
	"testing"

	"github.com/scorredoira/email"
	"github.com/zhnxin/emailagent"
)

func generateEmail() (*email.Message, error) {
	// compose the message
	m := email.NewMessage("Hi", "this is the body")
	// This is cantained in the message, which has no effect on email sending.
	m.From = mail.Address{Name: "nic name", Address: "username@aliyun.com"}
	m.To = []string{"target@aliyun.com"}

	// add attachments
	if err := m.Attach("agent.go"); err != nil {
		return nil, err
	}

	// add headers
	m.AddHeader("X-CUSTOMER-id", "xxxxx")
	return m, nil
}

func TestSSL(t *testing.T) {
	msg, err := generateEmail()
	if err != nil {
		t.Fatal(err)
	}
	agent := emailagent.New("exmaple@aliyun.com", "password", "smtp.aliyun.com", 465, true)
	//agent := emailagent.NewWithIdentify("identify","exmaple@aliyun.com", "password", "smtp.aliyun.com", 465, true)
	if err = agent.SendEmail(msg); err != nil {
		t.Fatal(err)
	}

}

func TestPlainAuth(t *testing.T) {
	msg, err := generateEmail()
	if err != nil {
		t.Fatal(err)
	}
	agent := emailagent.New("exmaple@aliyun.com", "password", "smtp.aliyun.com", 25, false)
	//agent := emailagent.NewWithIdentify("identify","exmaple@aliyun.com", "password", "smtp.aliyun.com", 25, false)
	if err = agent.SendEmail(msg); err != nil {
		t.Fatal(err)
	}
}


```