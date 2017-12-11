package adapter

import (
	"errors"
	"fmt"
	"strings"

	"github.com/pivotal-cf/on-demand-services-sdk/bosh"
	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
)

const (
	onlyStemcellAliasForService  = "service-only-stemcell"
	serviceJobIndex              = 0
	//syslogConfiguratorJobIndex = 1
	servicePropertyName          = "redis"
	propertyPass                 = "requirepass"
	propertyPort                 = "port"
	propertyConfigCommand        = "config_command"
)

var (

	servicePlanProperties        = []string{propertyConfigCommand, propertyPort}
	instanceGroupJobMap          = map[string][]string{
		                             "redis-instance":     []string{"redis"}, //instance-group-name : job-name
	                               }
)
//ManifestGenerator implements serviceadapter.ManifestGenerator
type ManifestGenerator struct{}

func (mg ManifestGenerator) GenerateManifest(
	deployment serviceadapter.ServiceDeployment,
	plan serviceadapter.Plan,
	params serviceadapter.RequestParameters,
	previousManifest *bosh.BoshManifest,
	previousPlan *serviceadapter.Plan,
) (bosh.BoshManifest, error) {
	generator := newServiceManifestGenerator(
		deployment,
		plan,
		params,
		previousManifest,
		previousPlan,
	)
	return generator.Generate()
}

type serviceManifestGenerator struct {
	deployment       serviceadapter.ServiceDeployment
	plan             serviceadapter.Plan
	params           serviceadapter.RequestParameters
	previousManifest *bosh.BoshManifest
	previousPlan     *serviceadapter.Plan
	planProperties   []string
}

//GenerateManifest creates a redis-service manifest
func (generator serviceManifestGenerator) Generate() (bosh.BoshManifest, error) {
	manifest, err := generator.getInitialManifest()
	if err != nil {
		return bosh.BoshManifest{}, err
	}

	manifest, err = generator.applyArbitraryParams(manifest)
	if err != nil {
		return bosh.BoshManifest{}, err
	}

	return generator.applyPlanConfig(manifest), nil
}

func (generator serviceManifestGenerator) getInitialManifest() (bosh.BoshManifest, error) {
	manifest, err := generator.newManifest()
	if err != nil {
		return bosh.BoshManifest{}, err
	}

	if generator.hasPreviousManifest() {
		manifest.InstanceGroups[0].Jobs[serviceJobIndex].Properties[servicePropertyName] = generator.previousManifest.InstanceGroups[0].Jobs[serviceJobIndex].Properties[servicePropertyName]
		manifest.InstanceGroups[0].AZs = generator.previousManifest.InstanceGroups[0].AZs
		manifest.InstanceGroups[0].Networks = generator.previousManifest.InstanceGroups[0].Networks
	}

	return manifest, nil
}

func (generator serviceManifestGenerator) newManifest() (bosh.BoshManifest, error) {
	instanceGroups, err := generator.generateInstanceGroupsBlock()
	if err != nil {
		return bosh.BoshManifest{}, err
	}

	return bosh.BoshManifest{
		Name:           generator.deployment.DeploymentName,
		Releases:       generator.generateReleasesBlock(),
		Stemcells:      generator.generateStemcellsBlock(),
		InstanceGroups: instanceGroups,
		Update:         generator.generateUpdateBlock(),
		Properties:     make(map[string]interface{}, 0),
	}, nil
}

func newServiceManifestGenerator(
	deployment serviceadapter.ServiceDeployment,
	plan serviceadapter.Plan,
	params serviceadapter.RequestParameters,
	previousManifest *bosh.BoshManifest,
	previousPlan *serviceadapter.Plan,
) serviceManifestGenerator {
	return serviceManifestGenerator{
		deployment:       deployment,
		plan:             plan,
		params:           params,
		previousManifest: previousManifest,
		previousPlan:     previousPlan,

		planProperties:   servicePlanProperties,
	}
}

func (generator serviceManifestGenerator) applyArbitraryParams(manifest bosh.BoshManifest) (bosh.BoshManifest, error) {
	/*if err := generator.validateArbitraryParams(); err != nil {
		return bosh.BoshManifest{}, err
	}*/

	for k, v := range generator.params.ArbitraryParams() {
		manifest.InstanceGroups[0].Jobs[serviceJobIndex].Properties[servicePropertyName].(map[interface{}]interface{})[k] = v
	}

	return manifest, nil
}

func (generator serviceManifestGenerator) applyPlanConfig(manifest bosh.BoshManifest) bosh.BoshManifest {
	if generator.plan.Properties == nil || generator.plan.Properties[servicePropertyName] == nil {
		return manifest
	}

	for _, property := range generator.planProperties {
		propertyValue := generator.plan.Properties[servicePropertyName].(map[string]interface{})[property]
		if propertyValue != nil {
			manifest.InstanceGroups[0].Jobs[serviceJobIndex].Properties[servicePropertyName].(map[interface{}]interface{})[property] = propertyValue
		}
	}
	manifest.InstanceGroups[0].PersistentDiskType = generator.plan.InstanceGroups[0].PersistentDiskType
	manifest.InstanceGroups[0].VMType = generator.plan.InstanceGroups[0].VMType

	return manifest
}

func (generator serviceManifestGenerator) generateStemcellsBlock() []bosh.Stemcell {
	return []bosh.Stemcell{{
		Alias:   onlyStemcellAliasForService,
		OS:      generator.deployment.Stemcell.OS,
		Version: generator.deployment.Stemcell.Version,
	}}
}

func (generator serviceManifestGenerator) generateReleasesBlock() []bosh.Release {
	var releases []bosh.Release
	for _, release := range generator.deployment.Releases {
		releases = append(releases, bosh.Release{Name: release.Name, Version: release.Version})
	}
	return releases
}

func (generator serviceManifestGenerator) generateInstanceGroupsBlock() ([]bosh.InstanceGroup, error) {
	password, err := newRequirepass()
	if err != nil {
		return nil, err
	}

	return generator.generateInstanceGroupsWithPassword(password)
}

func (generator serviceManifestGenerator) generateInstanceGroupsWithPassword(password string) ([]bosh.InstanceGroup, error) {
	instanceGroups, err := generator.generateInstanceGroupsWithNoProperties()
	if err != nil {
		return nil, err
	}

	instanceGroups[0].Jobs[serviceJobIndex].Properties = map[string]interface{}{
		servicePropertyName: map[interface{}]interface{}{propertyPass: password},
	}

	//if syslogProperties, ok := generator.plan.Properties["syslog_configurator"].(map[string]interface{}); ok {
	//	instanceGroups[0].Jobs[syslogConfiguratorJobIndex].Properties = map[string]interface{}{"syslog_configurator": syslogProperties}
	//}

	return instanceGroups, nil
}

func (generator serviceManifestGenerator) generateInstanceGroupsWithNoProperties() ([]bosh.InstanceGroup, error) {
	return serviceadapter.GenerateInstanceGroupsWithNoProperties(
		generator.plan.InstanceGroups,
		generator.deployment.Releases,
		onlyStemcellAliasForService,
		generator.mapInstanceGroupsToJobs(),
	)
}

func (generator serviceManifestGenerator) mapInstanceGroupsToJobs() map[string][]string {
	//jobs := append([]string{servicePropertyName}, generator.getJobsFromSyslogConfiguratorRelease()...)

	return instanceGroupJobMap
}

func (generator serviceManifestGenerator) generateUpdateBlock() bosh.Update {
	return bosh.Update{
		Canaries:        1,
		CanaryWatchTime: "30000-180000",
		UpdateWatchTime: "30000-180000",
		MaxInFlight:     6,
	}
}

func (generator serviceManifestGenerator) hasPreviousManifest() bool {
	return (generator.previousManifest != nil) && (generator.previousManifest.Name != "")
}

func (generator serviceManifestGenerator) getJobsFromSyslogConfiguratorRelease() []string {
	for _, r := range generator.deployment.Releases {
		if r.Name == "syslog-configurator" {
			return r.Jobs
		}
	}
	return nil
}

//------------------------------------------param check script start -----------------------------------------------

func (generator serviceManifestGenerator) validateArbitraryParams() error {
	err := generator.validateArbitraryParamsKeys()
	if err != nil {
		return err
	}
	return generator.validateArbitraryParamsValues()
}

func (generator serviceManifestGenerator) validateArbitraryParamsKeys() error {
	for k, _ := range generator.params.ArbitraryParams() {
		if err := validKey(k); err != nil {
			return err
		}
	}

	return nil
}

/**
param key check if you need
 */
func validKey(key string) error {
	allowedKeys := []string{
		"maxmemory-policy",
		"notify-keyspace-events",
		"lua-time-limit",
		"slowlog-max-len",
		"slowlog-log-slower-than",
	}

	for _, allowedKey := range allowedKeys {
		if allowedKey == key {
			return nil
		}
	}
	return fmt.Errorf("Invalid arbitrary param %q provided, must be one of the following: %v", key, allowedKeys)
}

/**
param value check if you need
 */
func (generator serviceManifestGenerator) validateArbitraryParamsValues() error {
	params := generator.params.ArbitraryParams()

	if err := validateMaxmemoryPolicy(params["maxmemory-policy"]); err != nil {
		return err
	}

	if err := validateNotifyKeyspace(params["notify-keyspace-events"]); err != nil {
		return err
	}

	intValidations := map[string]int{
		"lua-time-limit":          10000,
		"slowlog-log-slower-than": 20000,
	}

	for name, max := range intValidations {

		if err := validateIntLessThanMax(name, params[name], max); err != nil {
			return err
		}
	}

	if err := validateIntWithinRange("slowlog-max-len", params["slowlog-max-len"], 1, 2024); err != nil {
		return err
	}

	return nil
}

/**
param value check if you need
 */
func validateNotifyKeyspace(value interface{}) error {
	if value == nil {
		return nil
	}

	strval, ok := value.(string)
	if !ok {
		return errors.New("notify-keyspace-events must be of type string")
	}

	valid := "KEg$lshzxeA"
	for _, r := range strval {
		if !strings.ContainsRune(valid, r) {
			return fmt.Errorf("invalid flag %q specified for notify-keyspace-events, must be one of the following characters: %q", r, valid)
		}
	}

	return nil
}

/**
param value check if you need
 */
func validateIntWithinRange(name string, value interface{}, min, max int) error {
	if value == nil {
		return nil
	}

	fval, ok := value.(float64)
	if !ok {
		return fmt.Errorf("%s must be of type Number", name)
	}

	ival := int(fval)

	if ival < min || ival > max {
		return fmt.Errorf("out-of-range value for %s, must be within %d-%d", name, min, max)
	}

	return nil
}

/**
param value check if you need
 */
func validateIntLessThanMax(name string, value interface{}, max int) error {
	if value == nil {
		return nil
	}

	fval, ok := value.(float64)
	if !ok {
		return fmt.Errorf("%s must be of type Number", name)
	}

	ival := int(fval)

	if ival > max {
		return fmt.Errorf("out-of-range value for %s, must be less than %d", name, max)
	}

	return nil
}

/**
param value check if you need
 */
func validateMaxmemoryPolicy(value interface{}) error {
	if value == nil {
		return nil
	}

	maxmem, ok := value.(string)
	if !ok {
		return errors.New("maxmemory-policy must be of type string")
	}

	allowed := []string{
		"noeviction",
		"allkeys-lru",
		"volatile-lru",
		"allkeys-random",
		"volatile-ttl",
	}

	for _, a := range allowed {
		if strings.ToLower(maxmem) == a {
			return nil
		}
	}

	return fmt.Errorf("invalid value %q specified for maxmemory-policy", value)
}
//------------------------------------------param check script end -----------------------------------------------
