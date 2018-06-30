package email

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
)

// Timeout of connection dial
const Timeout = 10

// Credential stores smtp data
type Credential struct {
	Host  string
	Port  string
	Login string
	Pass  string
	TLS   bool
}

// Message contains data to prepare an email
type Message struct {
	From    mail.Address
	To      mail.Address
	Subject string
	Body    string
}

// Result is a type which is send to a channel
type Result struct {
	Client *smtp.Client
	Error  error
}

// NewClient creates a new mail client
func NewClient(credential Credential, c chan<- Result) {
	connString := net.JoinHostPort(credential.Host, credential.Port)
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         credential.Host,
	}

	// Connect to the remote SMTP server.
	client, err := smtp.Dial(connString)
	if err != nil {
		c <- Result{nil, err}
		return
	}

	// TLS
	if credential.TLS == true {
		err = client.StartTLS(tlsconfig)
		if err != nil {
			c <- Result{nil, err}
			return
		}
	}

	// Auth
	err = Auth(client, credential)
	if err != nil {
		c <- Result{nil, err}
		return
	}

	c <- Result{client, err}

}

// Auth plain auth
func Auth(c *smtp.Client, credential Credential) error {
	auth := smtp.PlainAuth("", credential.Login, credential.Pass, credential.Host)
	// Auth
	if err := c.Auth(auth); err != nil {
		return err
	}

	return nil
}

// prepareMessage format email message to send
func prepareMessage(msg Message) (string, error) {

	message := ""

	// Prepare message
	message += fmt.Sprintf("%s: %s\r\n", "From", msg.From.String())
	message += fmt.Sprintf("%s: %s\r\n", "To", msg.To.String())
	message += fmt.Sprintf("%s: %s\r\n", "Subject", msg.Subject)
	message += "\r\n" + msg.Body

	return message, nil
}

func validateEmail(e string) (bool, error) {
	if !strings.Contains(e, "@") {
		return false, fmt.Errorf("Invalid email address: %s", e)
	}

	host := strings.Split(e, "@")[1]

	_, err := net.LookupMX(host)

	if err != nil {
		return false, err
	}

	return true, nil
}

// Send sends an email
func Send(c *smtp.Client, m Message) error {
	// Validate emails
	if _, err := validateEmail(m.From.Address); err != nil {
		return err
	}

	if _, err := validateEmail(m.To.Address); err != nil {
		return err
	}

	// Set the sender and recipient
	if err := c.Mail(m.From.Address); err != nil {
		return err
	}

	if err := c.Rcpt(m.To.Address); err != nil {
		return err
	}

	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		return err
	}

	// Send message
	message, err := prepareMessage(m)
	if err != nil {
		return err
	}

	_, err = wc.Write([]byte(message))

	if err != nil {
		return err
	}

	err = wc.Close()
	if err != nil {
		return err
	}

	// Send the QUIT command and close the connection.
	err = c.Quit()
	if err != nil {
		return err
	}

	return nil
}
