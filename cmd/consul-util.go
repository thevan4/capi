package run

import consulapi "github.com/hashicorp/consul/api"

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
