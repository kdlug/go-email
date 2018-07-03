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

	"github.com/hashicorp/consul/api"
	"github.com/kdlug/email/config"
	"github.com/kdlug/email/email"
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
	// Get a new consul client
	consul, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		Error.Println(err)
	}

	// Get a handle to the KV API
	kv := consul.KV()

	// Lookup the pair
	// timeoutConfig, _, err := kv.Get(ServiceName+"/timeout", nil)
	// timeout, _ := strconv.Atoi(string(timeoutConfig.Value))
	// if err != nil {
	// 	Error.Println(err)
	// }
	timeout, err := config.GetTimeout(kv)
	if err != nil {
		Error.Println(err)
	}

	sender, err := config.GetSender(kv)
	if err != nil {
		Error.Println(err)
	}

	receipients, err := config.GetReceipients(kv)
	if err != nil {
		Error.Println(err)
	}
	fmt.Println(receipients)

	credentials, err := config.GetCredentials(kv)
	if err != nil {
		Error.Println(err)
	}

	fmt.Println(timeout)
	fmt.Println(sender)
	fmt.Println(credentials)
	return
	// fmt.Println("CONSUL_TIMEOUT:", os.Getenv("CONSUL_TIMEOUT"))
	// fmt.Println("CONSUL_ADDRESS:", os.Getenv("CONSUL_ADDRESS"))
	// fmt.Println("CONSUL_SCHEME:", os.Getenv("CONSUL_SCHEME"))

	msg := email.Message{
		mail.Address{"", sender},
		mail.Address{"", receipients[0]},
		"Test",
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
