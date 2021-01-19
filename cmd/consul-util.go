package run

import (
	"time"

	consulapi "github.com/hashicorp/consul/api"
)

func newConsulKVClient(consulAddress string, consulTimeout time.Duration) (*consulapi.KV, error) {
	// Get a new client
	clientAPI := consulapi.DefaultConfig()
	clientAPI.Address = consulAddress
	clientAPI.WaitTime = consulTimeout
	client, err := consulapi.NewClient(clientAPI)
	if err != nil {
		return nil, err
	}
	return client.KV(), nil
}

func isKeyExist(consulKVClient *consulapi.KV, keyPath string) (bool, error) {
	value, meta, err := consulKVClient.Get(keyPath, nil)
	if err != nil {
		return false, err
	}
	if value != nil && meta != nil {
		return true, nil
	}
	return false, nil
}
