package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/mail"
	"os"
	"time"

	"github.com/kdlug/email/email"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func Init(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	// Logger
	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	// fmt.Println("CONSUL_TIMEOUT:", os.Getenv("CONSUL_TIMEOUT"))
	// fmt.Println("CONSUL_ADDRESS:", os.Getenv("CONSUL_ADDRESS"))
	// fmt.Println("CONSUL_SCHEME:", os.Getenv("CONSUL_SCHEME"))

	msg := email.Message{
		mail.Address{"", "john.doe@gmail.com"},
		mail.Address{"", "recipient@gmail.com"},
		"Test Subject",
		"Test body",
	}

	c1 := email.Credential{"in-v3.mailjet.com", "2525", "username", "password", true}

	credentials := []email.Credential{}

	// credentials = append(credentials, c1, c2, c3)
	credentials = append(credentials, c1)

	channel := make(chan email.Result)

	for _, credential := range credentials {
		go email.NewClient(credential, channel)
	}

	// check all channels, set timeout
	for _, item := range credentials {
		select {
		case result := <-channel:

			if result.Error != nil {
				Error.Println(item.Host, result.Error)
			} else {
				if err := email.Send(result.Client, msg); err != nil {
					Error.Println(item.Host, ":", err)
				} else {
					Info.Println(item.Host, ":", "Email has been sent successfuly.")
				}

			}
		case <-time.After(email.Timeout * time.Second):
			fmt.Println(item.Host, ":", "Timeout")
		}
	}

}
