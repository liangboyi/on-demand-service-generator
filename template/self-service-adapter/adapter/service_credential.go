package adapter

import (
	"errors"
	"fmt"

	"github.com/pivotal-cf/on-demand-services-sdk/bosh"
)

type credential struct {
	cluster []map[string]interface{}
}


var invalidManifestErr = errors.New("On demand broker did not return a valid manifest.")

var hosts []string

func newCredential(deploymentTopology bosh.BoshVMs, manifest bosh.BoshManifest) (credential, error) {
	cred := credential{}

	var err error

	var port int
	port, err = getPort(manifest)
	if err != nil {
		return credential{}, err
	}

	hosts,err = getHosts(manifest, deploymentTopology)
	if err != nil {
		return credential{}, err
	}

	cred.cluster = make([]map[string]interface{},0)
	for _, host := range hosts {
		noPassMap := map[string]interface{}{
			"host": host,
			"port": port,
		}

		cred.cluster = append(cred.cluster,noPassMap)
	}

	return cred, nil
}

func (cred credential) bindingAdapter() map[string]interface{} {
	return map[string]interface{}{
		"cluster":     cred.cluster,
	}
}

func getHosts(manifest bosh.BoshManifest, deploymentTopology bosh.BoshVMs) ([]string, error) {
	if manifest.InstanceGroups == nil {
		return []string{}, invalidManifestErr
	}
	instanceGroupName := manifest.InstanceGroups[0].Name
	host := deploymentTopology[instanceGroupName]

	if host == nil {
		return []string{}, fmt.Errorf("On demand broker did not return instance groups matching %s, unable to find host.", instanceGroupName)
	}
	return host, nil
}

func getPassword(manifest bosh.BoshManifest) (string, error) {
	result,err := getProperty(manifest,servicePropertyName,propertyPort)
	return result.(string),err
}

func getPort(manifest bosh.BoshManifest) (int, error) {
	result,err := getProperty(manifest,servicePropertyName,propertyPort)
	return result.(int),err
}

func getProperty(manifest bosh.BoshManifest, parentNodeName string, property string) (interface{},error){
	if manifest.InstanceGroups == nil || manifest.InstanceGroups[0].Jobs == nil {
		return nil, invalidManifestErr
	}

	properties := getJobProperties(manifest)
	if properties == nil {
		return nil, errors.New("Missing "+parentNodeName+" properties.")
	}

	if properties == nil || getFromMap(properties, property) == nil {
		return nil, errors.New("Missing "+ property + " property.")
	}

	return getFromMap(properties, property), nil

}

func getJobProperties(manifest bosh.BoshManifest) interface{} {
	return manifest.InstanceGroups[0].Jobs[serviceJobIndex].Properties[servicePropertyName]
}

func getFromMap(m interface{}, key string) interface{} {
	return m.(map[interface{}]interface{})[key]
}
