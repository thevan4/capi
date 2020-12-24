package run

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	consulapi "github.com/hashicorp/consul/api"
)

func putRootPathToConsul(consulKVClient *consulapi.KV, rootPath string) error {
	isRootPathExist, err := isKeyExist(consulKVClient, rootPath)
	if err != nil {
		return err
	}
	if !isRootPathExist {
		proot := &consulapi.KVPair{Key: rootPath, Value: nil}
		_, err = consulKVClient.Put(proot, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func getServicesPaths(servicesDirectory string) ([]string, error) {
	services, err := ioutil.ReadDir(servicesDirectory)
	if err != nil {
		return nil, fmt.Errorf("get read dir %v, error: %v",
			servicesDirectory,
			err)
	}

	servicesPaths := []string{}
	for _, s := range services {
		if s.IsDir() {
			servicesPaths = append(servicesPaths, s.Name())
		}
	}
	return servicesPaths, nil
}

func getManifestsAndAppServers(servicesDirectory, servicePath string) ([]os.FileInfo, error) {
	manAndApps, err := ioutil.ReadDir(servicesDirectory + "/" + servicePath)
	if err != nil {
		return nil, fmt.Errorf("can't read dir %v, error: %v",
			servicesDirectory+"/"+servicePath,
			err)
	}
	return manAndApps, nil
}

func readManifest(servicesDirectory string,
	servicePath string,
	manifestName string) ([]byte, error) {
	dat, err := ioutil.ReadFile(servicesDirectory + "/" + servicePath + "/" + manifestName + ".json")
	if err != nil {
		return nil, fmt.Errorf("can't read file %v, got error: %v",
			servicesDirectory+"/"+servicePath+"/"+manifestName+".json",
			err)
	}
	return dat, nil
}

func putManifestToConsul(consulKVClient *consulapi.KV,
	servicePath string,
	rootPath string,
	manifestName string,
	fileData []byte,
	isForceUpdate bool) error {
	isServicePathExist, err := isKeyExist(consulKVClient, rootPath)
	if err != nil {
		return err
	}
	if !isServicePathExist {
		pServiceF := &consulapi.KVPair{Key: rootPath + servicePath + "/", Value: nil}
		_, err = consulKVClient.Put(pServiceF, nil)
		if err != nil {
			return fmt.Errorf("can't put service folder %v to consul, got error: %v",
				rootPath+servicePath+"/",
				err)
		}
	}

	isServiceManifestExist, err := isKeyExist(consulKVClient, rootPath+servicePath+"/"+manifestName)
	if err != nil {
		return err
	}

	if isForceUpdate || !isServiceManifestExist {
		pService := &consulapi.KVPair{Key: rootPath + servicePath + "/" + manifestName, Value: fileData}
		_, err = consulKVClient.Put(pService, nil)
		if err != nil {
			return fmt.Errorf("can't put service manifest %v to consul, got error: %v",
				rootPath+servicePath+"/"+manifestName,
				err)
		}
	}
	return nil
}

func readAndPutManifest(consulKVClient *consulapi.KV,
	servicesDirectory string,
	servicePath string,
	rootPath string,
	manifestName string,
	isForceUpdate bool) error {
	fileData, err := readManifest(servicesDirectory,
		servicePath,
		manifestName)
	if err != nil {
		return fmt.Errorf("read manifest error: %v", err)
	}

	if err := putManifestToConsul(consulKVClient,
		servicePath,
		rootPath,
		manifestName,
		fileData,
		isForceUpdate); err != nil {
		return fmt.Errorf("put manifest to consul error: %v",
			err)
	}

	return nil
}

func putConsulData(consulKVClient *consulapi.KV,
	servicesDirectory string,
	manifestName string,
	appSrvrsName string,
	rootPath string,
	isForceUpdate bool) error {
	servicesPaths, err := getServicesPaths(servicesDirectory)
	if err != nil {
		return err
	}

	for _, servicePath := range servicesPaths {
		manifestsAndAppServers, err := getManifestsAndAppServers(servicesDirectory, servicePath)
		if err != nil {
			return err
		}

		for _, maa := range manifestsAndAppServers {
			if maa.Name() == manifestName+".json" {
				if err := readAndPutManifest(consulKVClient,
					servicesDirectory,
					servicePath,
					rootPath,
					manifestName,
					isForceUpdate); err != nil {
					return fmt.Errorf("can't put manifest got error %v", err)
				}
			} else if maa.IsDir() && maa.Name() == appSrvrsName {
				// TODO: refactor
				isAppSrvrsPathExist, err := isKeyExist(consulKVClient, rootPath)
				if err != nil {
					return err
				}
				if !isAppSrvrsPathExist {
					appF := &consulapi.KVPair{Key: rootPath + servicePath + "/" + maa.Name() + "/", Value: nil}
					_, err = consulKVClient.Put(appF, nil)
					if err != nil {
						return fmt.Errorf("can't create app servers folder %v in consul: %v",
							rootPath+servicePath+"/"+maa.Name()+"/",
							err)
					}
				}

				appSrvsFiles, err := ioutil.ReadDir(servicesDirectory + "/" + servicePath + "/" + maa.Name())
				if err != nil {
					return fmt.Errorf("can't read app servers files: %v", err)
				}

				for _, app := range appSrvsFiles {
					dat, err := ioutil.ReadFile(servicesDirectory + "/" + servicePath + "/" + maa.Name() + "/" + app.Name())
					if err != nil {
						return fmt.Errorf("read app server file %v error: %v",
							servicesDirectory+"/"+servicePath+"/"+maa.Name()+"/"+app.Name(),
							err)
					}
					dApp := &ApplicationServerTransport{}
					if err := json.Unmarshal(dat, dApp); err != nil {
						return fmt.Errorf("can't unmarshal app server file %v data: %v",
							servicesDirectory+"/"+servicePath+"/"+maa.Name()+"/"+app.Name(),
							err)
					}

					isAppServerExist, err := isKeyExist(consulKVClient, rootPath+servicePath+"/"+maa.Name()+"/"+dApp.IP+":"+dApp.Port)
					if err != nil {
						return err
					}
					if isForceUpdate || !isAppServerExist {
						appI := &consulapi.KVPair{Key: rootPath + servicePath + "/" + maa.Name() + "/" + dApp.IP + ":" + dApp.Port, Value: dat}
						_, err = consulKVClient.Put(appI, nil)
						if err != nil {
							return fmt.Errorf("can't create app server config in consul: %v",
								err)
						}
					}
				}
			}
		}
	}
	return nil
}
