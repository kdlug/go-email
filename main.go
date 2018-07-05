// go get github.com/hashicorp/consul/api
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/mail"
	"os"
	"time"

	"github.com/kdlug/go-email/config"
	"github.com/kdlug/go-email/email"
	"github.com/kdlug/go-email/randgen"
)

// ServiceName
const ServiceName = "email-service"

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

	c := config.NewConfig()

	timeout, err := c.GetTimeout()

	if err != nil {
		Error.Println(err)
	}

	sender, err := c.GetSender()
	if err != nil {
		Error.Println(err)
	}

	receipients, err := c.GetReceipients()
	if err != nil {
		Error.Println(err)
	}
	fmt.Println(receipients)

	credentials, err := c.GetCredentials()
	if err != nil {
		Error.Println(err)
	}

	fmt.Println(timeout)
	fmt.Println(sender)
	fmt.Println(credentials)

	// fmt.Println("CONSUL_TIMEOUT:", os.Getenv("CONSUL_TIMEOUT"))
	// fmt.Println("CONSUL_ADDRESS:", os.Getenv("CONSUL_ADDRESS"))
	// fmt.Println("CONSUL_SCHEME:", os.Getenv("CONSUL_SCHEME"))

	msg := email.Message{
		mail.Address{"", sender},
		mail.Address{"", receipients[0]},
		randgen.GenerateStr(10, "[TEST-", "]"),
		"Test body",
	}

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
		case <-time.After(time.Duration(timeout) * time.Second):
			fmt.Println(item.Host, ":", "Timeout")
		}
	}

}
