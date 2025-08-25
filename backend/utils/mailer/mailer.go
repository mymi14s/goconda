package mailer

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"

	"github.com/beego/beego/v2/server/web"
)

// SendEmail sends an email message to one or more recipients.
// The message may be plain text or HTML. If it contains a '<' character,
// we will set the Content-Type to text/html; otherwise text/plain.
func SendEmail(msg string, recipients []string) error {
	host := web.AppConfig.DefaultString("smtp::host", "")
	port := web.AppConfig.DefaultInt("smtp::port", 587)
	user := web.AppConfig.DefaultString("smtp::username", "")
	pass := web.AppConfig.DefaultString("smtp::password", "")
	from := web.AppConfig.DefaultString("smtp::from", user)
	subject := web.AppConfig.DefaultString("smtp::default_subject", "Notification")

	if host == "" || user == "" || pass == "" || from == "":
		return fmt.Errorf("smtp not configured (host/user/pass/from)")
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	hello := host
	if i := strings.Index(host, ":"); i >= 0 { hello = host[:i] }

	// Heuristic for HTML
	contentType := "text/plain; charset=UTF-8"
	if strings.Contains(msg, "<") {
		contentType = "text/html; charset=UTF-8"
	}

	headers := map[string]string{
		"From": from,
		"To": strings.Join(recipients, ", "),
		"Subject": subject,
		"MIME-Version": "1.0",
		"Content-Type": contentType,
	}
	var buf bytes.Buffer
	for k, v := range headers {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	buf.WriteString("\r\n")
	buf.WriteString(msg)

	auth := smtp.PlainAuth("", user, pass, host)
	// Try TLS first (STARTTLS)
	conn, err := tlsDial(addr)
	if err != nil {
		// fallback: plain dial then STARTTLS
		c, err2 := smtp.Dial(addr)
		if err2 != nil { return err }
		defer c.Close()
		if err = c.Hello(hello); err != nil { return err }
		if ok, _ := c.Extension("STARTTLS"); ok {
			if err = c.StartTLS(&tls.Config{ServerName: host}); err != nil { return err }
		}
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil { return err }
		}
		if err = c.Mail(from); err != nil { return err }
		for _, rcpt := range recipients {
			if err = c.Rcpt(rcpt); err != nil { return err }
		}
		w, err := c.Data()
		if err != nil { return err }
		if _, err = w.Write(buf.Bytes()); err != nil { return err }
		return w.Close()
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil { return err }
	defer client.Close()

	if err = client.Auth(auth); err != nil { return err }
	if err = client.Mail(from); err != nil { return err }
	for _, rcpt := range recipients {
		if err = client.Rcpt(rcpt); err != nil { return err }
	}
	w, err := client.Data()
	if err != nil { return err }
	if _, err = w.Write(buf.Bytes()); err != nil { return err }
	return w.Close()
}

// tlsDial makes a TLS connection to addr without doing a prior plaintext SMTP handshake.
func tlsDial(addr string) (*tls.Conn, error) {
	dialer := &net.Dialer{}
	return tls.DialWithDialer(dialer, "tcp", addr, &tls.Config{InsecureSkipVerify: false})
}
