package config

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/kdlug/go-email/email"
)

// ServiceName
const ServiceName = "email-service"

// CKV holds pointer to KV
type CKV struct {
	*api.KV
}

// NewConfig returns CKV struct
func NewConfig() *CKV {
	// Get a new consul client
	api, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		fmt.Println(err)
	}

	return &CKV{api.KV()}
}

// GetTimeout timeout in seconds
func (kv *CKV) GetTimeout() (int, error) {

	// Lookup the pair
	config, _, err := kv.Get(ServiceName+"/timeout", nil)
	value, _ := strconv.Atoi(string(config.Value))

	if err != nil {
		return 0, err
	}

	return value, nil
}

// GetSender returs sender email
func (kv *CKV) GetSender() (string, error) {
	config, _, err := kv.Get(ServiceName+"/sender", nil)
	value := string(config.Value)

	if err != nil {
		return "", err
	}

	_, err = email.ValidateEmail(value)
	if err != nil {
		return "", err
	}

	return value, nil
}

// GetReceipients emails
// If at least one email is invalid it returns err
func (kv *CKV) GetReceipients() ([]string, error) {
	config, _, err := kv.Get(ServiceName+"/receipients", nil)
	if err != nil {
		return nil, err
	}

	values := strings.Split(string(config.Value), ",")

	for _, value := range values {
		_, err = email.ValidateEmail(value)
		if err != nil {
			return nil, err
		}
	}

	return values, nil
}

// GetCredentials gets credentials for accounts
func (kv *CKV) GetCredentials() ([]email.Credential, error) {
	config, _, err := kv.Get(ServiceName+"/credentials", nil)

	if err != nil {
		return nil, err
	}

	var credentials = []email.Credential{}

	jsonByte := []byte(string(config.Value))

	if err := json.Unmarshal(jsonByte, &credentials); err != nil {
		return nil, err
	}

	return credentials, nil
}
