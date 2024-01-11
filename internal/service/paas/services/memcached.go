package services

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
		loggingEnabled:     true,
		monitoringEnabled:  true,
	},
}
