package services

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type memcachedManager struct {
	service
}

var Memcached = memcachedManager{
	service{
		name:               ServiceTypeMemcached,
		class:              []string{ServiceClassCacher},
		defaultClass:       ServiceClassCacher,
		allowArbitrator:    false,
		allowBackup:        false,
		dataVolumeRequired: true,
		usersEnabled:       false,
		databasesEnabled:   false,
	},
}

func (s memcachedManager) serviceParametersSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"monitoring": {
			Type:     schema.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},
	}
}

func (s memcachedManager) serviceParametersDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"monitoring": {
			Type:     schema.TypeBool,
			Computed: true,
		},
	}
}

func (s memcachedManager) ExpandServiceParameters(tfMap map[string]interface{}) ServiceParameters {
	if tfMap == nil {
		return nil
	}

	serviceParameters := ServiceParameters{}

	if v, ok := tfMap["monitoring"].(bool); ok {
		serviceParameters["monitoring"] = v
	}

	return serviceParameters
}

func (s memcachedManager) flattenServiceParameters(serviceParameters ServiceParameters) map[string]interface{} {
	if serviceParameters == nil {
		return map[string]interface{}{}
	}

	tfMap := map[string]interface{}{}

	if v, ok := serviceParameters["monitoring"].(bool); ok {
		tfMap["monitoring"] = v
	}

	return tfMap
}
