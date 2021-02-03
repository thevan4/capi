package run

import (
	consulapi "github.com/hashicorp/consul/api"
)

func deleteConsulData(consulKVClient *consulapi.KV,
	delPath string,
	isForceUpdate bool) error {
	if isForceUpdate {
		_, err := consulKVClient.DeleteTree(delPath, nil)
		return err
	}
	_, err := consulKVClient.Delete(delPath, nil)
	return err
}
