---
subcategory: "PaaS"
layout: "aws"
page_title: "CROC Cloud: aws_paas_service"
description: |-
  Provides information about a PaaS service.
---

# Data Source: aws_paas_service

Provides information about a PaaS service.

## Example Usage

```terraform
data "aws_paas_service" "selected" {
  id = "fm-cluster-12345678"
}

output "paas_service_name" {
  value = data.aws_paas_service.selected.name
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The ID of the PaaS service (e.g. `fm-cluster-12345678`).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `auto_created_security_group_ids` - List of security group IDs that CROC Cloud created for the service.
* `backup_settings` - The backup settings for the service. The structure of this block is [described below](#backup_settings).
* `data_volume` - The data volume parameters for the service. The structure of this block is [described below](#data_volume).
* `endpoints` - List of endpoints for connecting to the service.
* `error_code` - The service error code.
* `error_description` - The detailed description of the service error.
* `high_availability` - Indicates whether this is a high-availability service.
* `instances` - List of instances that refers to the service. The structure of this block is [described below](#instances).
* `instance_type` - The instance type.
* `name` - The service name.
* `network_interface_ids` - List of network interface IDs.
* `root_volume` - The root volume parameters for the service. The structure of this block is [described below](#root_volume).
* `security_group_ids` - List of security group IDs that were specified for the service.
* `service_class` - The service class.
* `service_type` - The service type. This value determines which service parameters are included in the corresponding block.
    * `elasticsearch` - Elasticsearch parameters. The structure of this block is [described below](#elasticsearch-attribute-reference).
    * `memcached` - Memcached parameters. The structure of this block is [described below](#memcached-attribute-reference).
    * `pgsql` - PostgreSQL parameters. The structure of this block is [described below](#postgresql-attribute-reference).
    * `redis` - Redis parameters. The structure of this block is [described below](#redis-attribute-reference).
* `ssh_key_name` - The name of the SSH key for accessing instances.
* `status` - The current status of the service.
* `subnet_ids` - List of subnet IDs.
* `supported_features` - List of service features.
* `total_cpu_count` - Total number of CPU cores in use.
* `total_memory` - Total RAM in use in MiB.
* `vpc_id` - The ID of the VPC.

### backup_settings

The `backup_settings` block has the following structure:

* `bucket_name` - The name of the bucket in object storage where the service backup is saved.
* `enabled` -  Indicates whether backup is enabled for the service.
* `expiration_days` - The backup retention period in days.
* `notification_email` - The email address to which a notification that backup was created is sent.
* `start_time` - The time when the daily backup process starts. It is set as a string in the HH:MM format Moscow time.
* `user_id` - The ID of a user with write permissions to the bucket in object storage.
* `user_login` - The login of a user with write permissions to the bucket in object storage.

### data_volume

The `data_volume` block has the following structure:

* `iops` - The number of read/write operations per second for the data volume.
* `size` - The size of the data volume in GiB.
* `type` - The type of the data volume.

### instances

The `instances` block has the following structure:

* `endpoint` - The service endpoint on the instance.
* `index` - The instance index.
* `instance_id` - The ID of the instance.
* `interface_id` - The ID of the instance network interface.
* `name` - The instance name.
* `private_ip` - The private IP address of the instance.
* `role` - The instance role.
* `status` - The current status of the instance.

### root_volume

The `root_volume` block has the following structure:

* `iops` - The number of read/write operations per second for the root volume.
* `size` - The size of the root volume in GiB.
* `type` - The type of the root volume.

## Elasticsearch Attribute Reference

~> **Note** The following attributes contain default parameter values or user-defined values when the service is created.

In addition to the common attributes for all services [described above](#attribute-reference),
the following attributes are exported only for an Elasticsearch service:

* `class` - The service class.
* `kibana` - Indicates whether the Kibana deployment is enabled.
* `logging` - The logging settings for the service. The structure of this block is [described below](#logging).
* `monitoring` - The monitoring settings for the service. The structure of this block is [described below](#monitoring).
* `options` - Other Elasticsearch parameters.
* `password` - The Elasticsearch user password.
* `version` - The installed version.

## Memcached Attribute Reference

~> **Note** The following attributes contain default parameter values or user-defined values when the service is created.

In addition to the common attributes for all services [described above](#attribute-reference),
the following attributes are exported only for a Memcached service:

* `class` - The service class.
* `logging` - The logging settings for the service. The structure of this block is [described below](#logging).
* `monitoring` - The monitoring settings for the service. The structure of this block is [described below](#monitoring).

## PostgreSQL Attribute Reference

~> **Note** The following attributes contain default parameter values or user-defined values when the service is created.

In addition to the common attributes for all services [described above](#attribute-reference),
the following attributes are exported only for a PostgreSQL service:

* `autovacuum` - Indicates whether the server must run the autovacuum launcher daemon.
* `autovacuum_max_workers` - The maximum number of autovacuum processes (other than the autovacuum launcher)
  that can run simultaneously.
* `autovacuum_vacuum_cost_delay` - The cost delay value in milliseconds used in automatic `VACUUM` operations.
* `autovacuum_vacuum_cost_limit` - The cost limit value used in automatic `VACUUM` operations.
* `autovacuum_analyze_scale_factor` - The fraction of the table size to add to `autovacuum_analyze_threshold`
  when deciding whether to trigger an `ANALYZE`.
* `autovacuum_vacuum_scale_factor` - The fraction of the table size to add to `autovacuum_vacuum_threshold`
  when deciding whether to trigger a `VACUUM`.
* `class` - The service class.
* `database` - List of PostgreSQL databases with parameters. The structure of this block is [described below](#postgresql-database).
* `effective_cache_size` - The planner’s assumption about the effective size of the disk cache
  that is available to a single query.
* `effective_io_concurrency` -  The number of concurrent disk I/O operations.
* `logging` - The logging settings for the service. The structure of this block is [described below](#logging).
* `maintenance_work_mem` -  The maximum amount of memory in bytes used by maintenance operations,
  such as `VACUUM`, `CREATE INDEX`, and `ALTER TABLE ADD FOREIGN KEY`.
* `max_connections` - The maximum number of simultaneous connections to the database server.
* `max_wal_size` - The maximum size in bytes that WAL can reach at automatic checkpoints.
* `max_parallel_maintenance_workers` - The maximum number of parallel workers that a single utility command can start.
* `max_parallel_workers` - The maximum number of workers that the system can support for parallel operations.
* `max_parallel_workers_per_gather` - The maximum number of workers that a single _Gather_ node can start.
* `max_worker_processes` - The maximum number of background processes that the system can support.
* `min_wal_size` - The minimum size in bytes to shrink the WAL to. As long as WAL disk usage stays below this setting,
  old WAL files are always recycled for future use at a checkpoint, rather than removed.
* `monitoring` - The monitoring settings for the service. The structure of this block is [described below](#monitoring).
* `options` - Other PostgreSQL parameters.
* `replication_mode` - The replication mode in the _Patroni_ cluster.
* `shared_buffers` - The amount of memory the database server uses for shared memory buffers.
* `user` - List of PostgreSQL users with parameters. The structure of this block is [described below](#postgresql-user).
* `version` - The installed version.
* `wal_buffers` - The amount of shared memory used for WAL data not yet written to a volume.
* `wal_keep_segments` - The minimum number of log files segments that must be kept in the _pg_xlog_ directory,
  in case a standby server needs to fetch them for streaming replication.
* `work_mem` - The base maximum amount of memory in bytes to be used by a query operation (such as a sort or hash table)
  before writing to temporary disk files.

### PostgreSQL database

The `database` block has the following structure:

* `backup_enabled` - Indicates whether backup is enabled for the database.
* `backup_id` - The database backup ID.
* `backup_db_name` - The name of a database from the backup specified in the `backup_id` parameter.
* `encoding` - The database encoding.
* `extensions` - List of extensions for the database.
* `id` - The ID of the database.
* `locale` - The database locale.
* `name` - The database name.
* `owner` - The name of the user who is the database owner.
* `user` - List of database users with parameters. The structure of this block is [described below](#postgresql-database-user).

### PostgreSQL database user

The `user` block has the following structure:

* `id` - The ID of the user.
* `name` - The PostgreSQL user name.

### PostgreSQL user

The `user` block has the following structure:

* `id` - The ID of the user.
* `name` - The PostgreSQL user name.
* `password` - The PostgreSQL user password.

## Redis Attribute Reference

~> **Note** The following attributes contain default parameter values or user-defined values when the service is created.

In addition to the common attributes for all services [described above](#attribute-reference),
the following attributes are exported only for a Redis service:

* `class` - The service class.
* `cluster_type` - The clustering option.
* `databases` - The number of databases.
* `logging` - The logging settings for the service. The structure of this block is [described below](#logging).
* `maxmemory_policy` - The memory management mode.
* `monitoring` - The monitoring settings for the service. The structure of this block is [described below](#monitoring).
* `options` - Other Redis parameters.
* `password` - The Redis user password.
* `persistence_aof` - Indicates whether the AOF storage mode is enabled.
* `persistence_rdb` - Indicates whether the RDB storage mode is enabled.
* `timeout` - The time in seconds during which the connection to an inactive client is retained.
* `tcp_backlog` - The size of a connection queue.
* `tcp_keepalive` - The time in seconds during which the service sends ACKs to detect dead peers (unreachable clients).
* `version` - The installed version.

## Common Service Attribute Reference

### logging

The `logging` block has the following structure:

* `log_to` - The ID of the logging service.
* `logging_tags` - List of tags that are assigned to the log records of the service.

### monitoring

The `monitoring` block has the following structure:

* `monitor_by` - The ID of the monitoring service.
* `monitoring_labels` - Map containing labels that are assigned to the metrics of the service.
