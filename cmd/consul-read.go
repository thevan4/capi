package run

import (
	"encoding/json"
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
)

func readConsulData(consulKVClient *consulapi.KV,
	rootPath string,
	rawConsulAppServersFolderName string,
	serviceManifest string) ([]*ServiceTransport, error) {
	appServersPath := rawConsulAppServersFolderName + "/"

	balancingServices, meta, err := consulKVClient.Keys(rootPath, "/", nil)
	if err != nil {
		return nil, fmt.Errorf("can't get consul services by kv path %v, error: %v",
			rootPath,
			err)
	}
	if balancingServices == nil || meta == nil || len(balancingServices) <= 1 {
		return nil, nil
	}

	balancingServicesTransportArray := make([]*ServiceTransport, 0, len(balancingServices)-1)
	for _, bsPath := range balancingServices {
		if bsPath == rootPath {
			continue
		}

		applicationServersPaths, _, err := consulKVClient.Keys(bsPath+appServersPath, "/", nil)
		if err != nil {
			return nil, fmt.Errorf("can't get consul application servers paths by kv path %v, error: %v",
				bsPath+appServersPath,
				err)
		}
		if len(applicationServersPaths) == 0 {
			return nil, fmt.Errorf("consul application servers paths by kv path %v is empty",
				bsPath+appServersPath)
		}
		applicationServersTransportArray := make([]*ApplicationServerTransport, 0, len(applicationServersPaths)-1)
		for _, applicationServersPath := range applicationServersPaths {
			if applicationServersPath == bsPath+appServersPath {
				continue
			}
			applicationServerPair, _, err := consulKVClient.Get(applicationServersPath, nil)
			if err != nil {
				return nil, fmt.Errorf("can't get consul application servers pair by kv path %v, error: %v",
					applicationServersPath,
					err)
			}
			applicationServerTransport := &ApplicationServerTransport{}
			if err := json.Unmarshal(applicationServerPair.Value, applicationServerTransport); err != nil {
				return nil, fmt.Errorf("can't unmarshal consul data by kv path %v, error: %v",
					applicationServersPath,
					err)
			}
			applicationServersTransportArray = append(applicationServersTransportArray, applicationServerTransport)
		}

		serviceManifestPair, _, err := consulKVClient.Get(bsPath+serviceManifest, nil)
		if err != nil {
			return nil, fmt.Errorf("can't service manifest pair by kv path %v, error: %v",
				bsPath+serviceManifest,
				err)
		}
		balancingServiceTransport := &ServiceTransport{}
		if err := json.Unmarshal(serviceManifestPair.Value, balancingServiceTransport); err != nil {
			return nil, fmt.Errorf("can't unmarshal consul data by kv path %v, error: %v",
				bsPath+serviceManifest,
				err)
		}
		balancingServiceTransport.ApplicationServersTransport = applicationServersTransportArray
		balancingServicesTransportArray = append(balancingServicesTransportArray, balancingServiceTransport)
	}
	return balancingServicesTransportArray, nil
}
