package config

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/kdlug/email/email"
)

// ServiceName
const ServiceName = "email-service"

type Handle api.KV

// GetTimeout timeout in seconds
func GetTimeout(kv *api.KV) (int, error) {

	// Lookup the pair
	config, _, err := kv.Get(ServiceName+"/timeout", nil)
	value, _ := strconv.Atoi(string(config.Value))

	if err != nil {
		return 0, err
	}

	return value, nil
}

// GetSender returs sender email
func GetSender(kv *api.KV) (string, error) {
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
func GetReceipients(kv *api.KV) ([]string, error) {
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

// GetCredentials get credentials for accounts
func GetCredentials(kv *api.KV) ([]email.Credential, error) {
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
