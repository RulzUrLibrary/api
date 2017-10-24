package smtp

import (
	"bytes"
	"fmt"
	"github.com/labstack/echo"
	"net/mail"
	"net/smtp"
	"strings"
)

const alt_tld = "contact@rulz.bar"

type Configuration struct {
	User     string
	Password string
	Host     string
	Port     int
}

type Smtp struct {
	Auth   smtp.Auth
	Logger echo.Logger
	conf   Configuration
}

type Addresses []mail.Address

func (a Addresses) Strings() (addresses []string) {
	for _, address := range a {
		addresses = append(addresses, address.Address)
	}
	return
}

func (a Addresses) String() string {
	addresses := []string{}
	for _, address := range a {
		addresses = append(addresses, address.String())
	}
	return strings.Join(addresses, ", ")
}

type Mail struct {
	to      Addresses
	from    mail.Address
	subject string
	body    []byte
	*Smtp
}

func (s *Smtp) NewMail(from string) *Mail {
	//return &Mail{from: mail.Address{from, s.conf.User}, Smtp: s}
	return &Mail{from: mail.Address{from, alt_tld /* should be replaced by s.User*/}, Smtp: s}
}

func (s *Smtp) SendMail(to []string, msg []byte) error {
	s.Logger.Debug(msg)
	host := fmt.Sprintf("%s:%d", s.conf.Host, s.conf.Port)
	return smtp.SendMail(host, s.Auth, s.conf.User, to, msg)
}

func (m *Mail) To(name, address string) *Mail {
	m.to = append(m.to, mail.Address{name, address})
	return m
}

func (m *Mail) Subject(object string) *Mail {
	addr := mail.Address{object, ""}
	m.subject = strings.Trim(addr.String(), " <>@")
	return m
}

func (m *Mail) Body(content []byte) *Mail {
	m.body = content
	return m
}

func (m *Mail) Body(content []string) *Mail {
	m.body = content
	return m
}

func (m *Mail) Send() error {
	header := map[string]string{}
	msg := ""

	header["From"] = m.from.String()
	header["To"] = m.to.String()
	header["Subject"] = m.subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""

	for k, v := range header {
		msg += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	body := bytes.NewBufferString(msg + "\r\n")
	body.Write(m.body)

	return m.SendMail(m.to.Strings(), body.Bytes())
}

func New(l echo.Logger, c Configuration) *Smtp {
	auth := smtp.PlainAuth("", c.User, c.Password, c.Host)
	return &Smtp{auth, l, c}
}
