package mailer

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"strings"

	"github.com/beego/beego/v2/server/web"
	"github.com/mymi14s/goconda/utils"
)

var ERROR_TITLE = "Email Sender"

// helper: log in one place and return the same error
func logErr(err error, context string) error {
	if err == nil {
		return nil
	}
	utils.LogError(map[string]any{
		"title":   ERROR_TITLE,
		"error":   err,
		"context": context,
	})
	return err
}

// SendEmail kicks off an async email send and returns immediately (fire-and-forget).
// Any errors are logged inside the goroutine.
func SendEmail(msg string, recipients []string, subject string) error {
	if subject == "" {
		subject = web.AppConfig.DefaultString("smtp::default_subject", "Notification")
	}
	go func(m string, rcpts []string, subj string) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("mailer: recovered from panic in SendEmail: %v", r)
			}
		}()
		if err := sendEmailSync(m, rcpts, subj); err != nil {
			log.Printf("mailer: failed to send email: %v", err)
		}
	}(msg, recipients, subject)
	return nil
}

// sendEmailSync contains the original sending logic (blocking).
func sendEmailSync(msg string, recipients []string, subject string) error {
	host := web.AppConfig.DefaultString("smtp::host", "")
	port := web.AppConfig.DefaultInt("smtp::port", 587)
	user := web.AppConfig.DefaultString("smtp::username", "")
	pass := web.AppConfig.DefaultString("smtp::password", "")
	from := web.AppConfig.DefaultString("smtp::from", user)

	if host == "" || user == "" || pass == "" || from == "" {
		return logErr(fmt.Errorf("smtp not configured (host/user/pass/from)"), "config")
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	hello := host
	if i := strings.Index(host, ":"); i >= 0 {
		hello = host[:i]
	}

	// Basic heuristic for HTML
	contentType := "text/plain; charset=UTF-8"
	if strings.Contains(strings.ToLower(msg), "<html") || strings.Contains(msg, "<") {
		contentType = "text/html; charset=UTF-8"
	}

	// Build message
	var buf bytes.Buffer
	writeHeader := func(k, v string) { buf.WriteString(fmt.Sprintf("%s: %s\r\n", k, v)) }
	writeHeader("From", from)
	writeHeader("To", strings.Join(recipients, ", "))
	writeHeader("Subject", subject)
	writeHeader("MIME-Version", "1.0")
	writeHeader("Content-Type", contentType)
	buf.WriteString("\r\n")
	buf.WriteString(msg)

	auth := smtp.PlainAuth("", user, pass, host)

	// Try implicit TLS (e.g., 465) first
	conn, err := tlsDial(addr)
	if err != nil {
		// Fallback: plain dial then STARTTLS (e.g., 587)
		c, err2 := smtp.Dial(addr)
		if err2 != nil {
			// return the original TLS error for better signal; include context
			return logErr(err, "tlsDial+Dial")
		}
		defer c.Close()

		if err = c.Hello(hello); err != nil {
			return logErr(err, "hello")
		}
		if ok, _ := c.Extension("STARTTLS"); ok {
			if err = c.StartTLS(&tls.Config{ServerName: host}); err != nil {
				return logErr(err, "starttls")
			}
		}
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				return logErr(err, "auth(starttls path)")
			}
		}
		if err = c.Mail(from); err != nil {
			return logErr(err, "mail")
		}
		for _, rcpt := range recipients {
			if err = c.Rcpt(rcpt); err != nil {
				return logErr(err, "rcpt")
			}
		}
		w, err := c.Data()
		if err != nil {
			return logErr(err, "data")
		}
		if _, err = w.Write(buf.Bytes()); err != nil {
			return logErr(err, "write")
		}
		return logErr(w.Close(), "closeWriter")
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return logErr(err, "newClient")
	}
	defer client.Close()

	if err = client.Auth(auth); err != nil {
		return logErr(err, "auth(implicit tls)")
	}
	if err = client.Mail(from); err != nil {
		return logErr(err, "mail(implicit tls)")
	}
	for _, rcpt := range recipients {
		if err = client.Rcpt(rcpt); err != nil {
			return logErr(err, "rcpt(implicit tls)")
		}
	}
	w, err := client.Data()
	if err != nil {
		return logErr(err, "data(implicit tls)")
	}
	if _, err = w.Write(buf.Bytes()); err != nil {
		return logErr(err, "write(implicit tls)")
	}
	return logErr(w.Close(), "closeWriter(implicit tls)")
}

// tlsDial makes a TLS connection to addr without doing a prior plaintext SMTP handshake.
func tlsDial(addr string) (*tls.Conn, error) {
	dialer := &net.Dialer{}
	return tls.DialWithDialer(dialer, "tcp", addr, &tls.Config{InsecureSkipVerify: false})
}
