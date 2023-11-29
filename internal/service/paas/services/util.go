package services

import (
	"log"
)

const (
	ServiceTypeElasticSearch = "elasticsearch"
	ServiceTypeMemcached     = "memcached"
	ServiceTypeMongoDB       = "mongodb"
	ServiceTypeMySQL         = "mysql"
	ServiceTypePostgreSQL    = "pgsql"
	ServiceTypeRabbitMQ      = "rabbitmq"
	ServiceTypeRedis         = "redis"
)

func ServiceTypeValues() []string {
	return []string{
		ServiceTypeElasticSearch,
		ServiceTypeMemcached,
		ServiceTypeMongoDB,
		ServiceTypeMySQL,
		ServiceTypePostgreSQL,
		ServiceTypeRabbitMQ,
		ServiceTypeRedis,
	}
}

const (
	ServiceClassCacher        = "cacher"
	ServiceClassDatabase      = "database"
	ServiceClassMessageBroker = "message_broker"
	ServiceClassSearch        = "search"
)

func ServiceClassValues() []string {
	return []string{
		ServiceClassCacher,
		ServiceClassDatabase,
		ServiceClassMessageBroker,
		ServiceClassSearch,
	}
}

const (
	Kilobyte = 1024
	Megabyte = 1024 * Kilobyte
	Gigabyte = 1024 * Megabyte
)

// Map with ServiceManager objects for each supported PaaS service.
var managers = map[string]ServiceManager{
	ElasticSearch.ServiceType(): ElasticSearch,
	Memcached.ServiceType():     Memcached,
	MongoDB.ServiceType():       MongoDB,
	MySQL.ServiceType():         MySQL,
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
