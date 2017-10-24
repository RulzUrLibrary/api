package smtp

import (
	"bytes"
	"fmt"
	"github.com/labstack/echo"
	"mime/quotedprintable"
	"net/mail"
	"net/smtp"
	"strings"
)

const alt_tld = "contact@rulz.bar"

func encodeRFC2047(String string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{String, ""}
	return strings.Trim(addr.String(), " <>@\"")
}

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
	host := fmt.Sprintf("%s:%d", s.conf.Host, s.conf.Port)
	return smtp.SendMail(host, s.Auth, s.conf.User, to, msg)
}

func (m *Mail) To(name, address string) *Mail {
	m.to = append(m.to, mail.Address{name, address})
	return m
}

func (m *Mail) Subject(object string) *Mail {
	m.subject = encodeRFC2047(object)
	return m
}

func (m *Mail) Body(content []byte) *Mail {
	m.body = content
	return m
}

func (m *Mail) BodyS(content string) *Mail {
	return m.Body([]byte(content))
}

func (m *Mail) Send() error {
	header := map[string]string{}
	msg := ""

	header["From"] = m.from.String()
	header["To"] = m.to.String()
	header["Subject"] = m.subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "quoted-printable"
	header["Content-Disposition"] = "inline"

	for k, v := range header {
		msg += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	body := bytes.NewBufferString(msg + "\r\n")

	encode := quotedprintable.NewWriter(body)
	encode.Write(m.body)
	encode.Close()

	return m.SendMail(m.to.Strings(), body.Bytes())
}

func New(l echo.Logger, c Configuration) *Smtp {
	auth := smtp.PlainAuth("", c.User, c.Password, c.Host)
	return &Smtp{auth, l, c}
}
