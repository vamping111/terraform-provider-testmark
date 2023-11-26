package services

import (
	"log"
)

const (
	ServiceTypeElasticSearch = "elasticsearch"
	ServiceTypeMemcached     = "memcached"
	ServiceTypePostgreSQL    = "pgsql"
	ServiceTypeRabbitMQ      = "rabbitmq"
	ServiceTypeRedis         = "redis"
)

const (
	ServiceClassCacher        = "cacher"
	ServiceClassDatabase      = "database"
	ServiceClassMessageBroker = "message_broker"
	ServiceClassSearch        = "search"
)

const (
	Kilobyte = 1024
	Megabyte = 1024 * Kilobyte
	Gigabyte = 1024 * Megabyte
)

// Map with ServiceManager objects for each supported PaaS service.
var managers = map[string]ServiceManager{
	ElasticSearch.ServiceType(): ElasticSearch,
	Memcached.ServiceType():     Memcached,
	PostgreSQL.ServiceType():    PostgreSQL,
	Redis.ServiceType():         Redis,
	RabbitMQ.ServiceType():      RabbitMQ,
}

func ManagedServiceTypes() []string {
	keys := make([]string, 0, len(managers))

	for k := range managers {
		keys = append(keys, k)
	}
	return keys
}

func Manager(serviceType string) ServiceManager {
	if v, ok := managers[serviceType]; ok {
		return v
	}

	log.Printf("[ERROR] Unknown service type: %s", serviceType)
	return nil
}
