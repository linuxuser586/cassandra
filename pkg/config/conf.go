package config

import (
	"github.com/linuxuser586/cassandra/pkg/config/commit"
	"github.com/linuxuser586/cassandra/pkg/config/disk"
)

// Compression is the YAML representation for CassandraYAML.
type Compression struct {
	ClassName  CompressionType `yaml:"class_name"`
	Parameters []string        `yaml:"parameters,omitempty"`
}

// SeedProvider contains the seed provider implementation class and seeds.
type SeedProvider struct {
	ClassName  string           `yaml:"class_name"`
	Parameters []SeedParameters `yaml:"parameters"`
}

// SeedParameters contains the seed IPs and/or hosts.
type SeedParameters struct {
	Seeds string `yaml:"seeds"`
}

// ServerEncryptionOptions hold all the data for Cassandra server encryption.
type ServerEncryptionOptions struct {
	InternodeEncryption         string   `yaml:"internode_encryption"`
	Keystore                    string   `yaml:"keystore"`
	KeystorePassword            string   `yaml:"keystore_password"`
	Truststore                  string   `yaml:"truststore"`
	TruststorePassword          string   `yaml:"truststore_password"`
	Protocol                    string   `yaml:"protocol,omitempty"`
	Algorithm                   string   `yaml:"algorithm,omitempty"`
	StoreType                   string   `yaml:"store_type,omitempty"`
	CipherSuites                []string `yaml:"cipher_suites,omitempty"`
	RequireClientAuth           bool     `yaml:"require_client_auth,omitempty"`
	RequireEndpointVerification bool     `yaml:"require_endpoint_verification,omitempty"`
}

// ClientEncryptionOptions hold all the data for Cassandra client encryption.
type ClientEncryptionOptions struct {
	Enabled            bool     `yaml:"enabled"`
	Optional           bool     `yaml:"optional"`
	Keystore           string   `yaml:"keystore"`
	KeystorePassword   string   `yaml:"keystore_password"`
	RequireClientAuth  bool     `yaml:"require_client_auth,omitempty"`
	Truststore         string   `yaml:"truststore,omitempty"`
	TruststorePassword string   `yaml:"truststore_password,omitempty"`
	Protocol           string   `yaml:"protocol,omitempty"`
	Algorithm          string   `yaml:"algorithm,omitempty"`
	StoreType          string   `yaml:"store_type,omitempty"`
	CipherSuites       []string `yaml:"cipher_suites,omitempty"`
}

// TransparentDataEncryptionOptions enables encrypting data at-rest (on disk).
type TransparentDataEncryptionOptions struct {
	Enabled       bool          `yaml:"enabled"`
	ChunkLengthKB int64         `yaml:"chunk_length_kb"`
	Cipher        string        `yaml:"cipher"`
	KeyAlias      string        `yaml:"key_alias"`
	KeyProvider   []KeyProvider `yaml:"key_provider"`
}

// KeyProvider contains the data for the TransparentDataEncryptionOptions
type KeyProvider struct {
	ClassName  string                  `yaml:"class_name"`
	Parameters []KeyProviderParameters `yaml:"parameters"`
}

// KeyProviderParameters are the KeyProvider parameters
type KeyProviderParameters struct {
	Keystore         string `yaml:"keystore"`
	KeystorePassword string `yaml:"keystore_password"`
	StoreType        string `yaml:"store_type"`
	KeyPassword      string `yaml:"key_password"`
}

// BackPressureStrategy is the back-pressure strategy applied.
type BackPressureStrategy struct {
	ClassName  string                   `yaml:"class_name"`
	Parameters []map[string]interface{} `yaml:"parameters"`
}

// RequestSchedulerOptions for the request scheduler
type RequestSchedulerOptions struct {
	ThrottleLimit int64            `yaml:"throttle_limit"`
	DefaultWeight int64            `yaml:"default_weight"`
	Weights       map[string]int64 `yaml:"weights"`
}

// CassandraYAML is the YAML file for Cassandra.
type CassandraYAML struct {
	ClusterName                                   string                           `yaml:"cluster_name"`
	NumTokens                                     int64                            `yaml:"num_tokens"`
	AllocateTokensForKeyspace                     string                           `yaml:"allocate_tokens_for_keyspace,omitempty"`
	InitialToken                                  string                           `yaml:"initial_token,omitempty"`
	HintedHandoffEnabled                          bool                             `yaml:"hinted_handoff_enabled"`
	HintedHandoffDisabledDatacenters              []string                         `yaml:"hinted_handoff_disabled_datacenters,omitempty"`
	MaxHintWindowInMS                             int64                            `yaml:"max_hint_window_in_ms"`
	HintedHandOffThrottleInKB                     int64                            `yaml:"hinted_handoff_throttle_in_kb"`
	MaxHintsDeliveryThreads                       int64                            `yaml:"max_hints_delivery_threads"`
	HintsDirectory                                string                           `yaml:"hints_directory,omitempty"`
	HintsFlushPeriodInMS                          int64                            `yaml:"hints_flush_period_in_ms"`
	MaxHintsFileSizeInMB                          int64                            `yaml:"max_hints_file_size_in_mb"`
	HintsCompression                              []Compression                    `yaml:"hints_compression,omitempty"`
	BatchlogReplayThrottleInKB                    int64                            `yaml:"batchlog_replay_throttle_in_kb"`
	Authenticator                                 string                           `yaml:"authenticator"`
	Authorizer                                    string                           `yaml:"authorizer"`
	RoleManager                                   string                           `yaml:"role_manager"`
	RolesValidityInMS                             int64                            `yaml:"roles_validity_in_ms"`
	RolesUpdateIntervalInMS                       int64                            `yaml:"roles_update_interval_in_ms,omitempty"`
	PermissionsValidityInMS                       int64                            `yaml:"permissions_validity_in_ms"`
	PermissionsUpdateIntervalInMS                 int64                            `yaml:"permissions_update_interval_in_ms,omitempty"`
	CredentialsValidityInMS                       int64                            `yaml:"credentials_validity_in_ms"`
	CredentialsUpdateIntervalInMS                 int64                            `yaml:"credentials_update_interval_in_ms,omitempty"`
	Partitioner                                   string                           `yaml:"partitioner"`
	DataFileDirectories                           []string                         `yaml:"data_file_directories,omitempty"`
	CommitlogDirectory                            string                           `yaml:"commitlog_directory,omitempty"`
	CDCEnabled                                    bool                             `yaml:"cdc_enabled"`
	CDCRawDirectory                               string                           `yaml:"cdc_raw_directory,omitempty"`
	DiskFailurePolicy                             disk.FailurePolicy               `yaml:"disk_failure_policy"`
	CommitFailurePolicy                           commit.FailurePolicy             `yaml:"commit_failure_policy"`
	PreparedStatementsCacheSizeMB                 int64                            `yaml:"prepared_statements_cache_size_mb,omitempty"`
	ThriftPreparedStatementsCacheSizeMB           int64                            `yaml:"thrift_prepared_statements_cache_size_mb,omitempty"`
	KeyCacheSizeInMB                              int64                            `yaml:"key_cache_size_in_mb,omitempty"`
	KeyCacheSavePeriod                            int64                            `yaml:"key_cache_save_period"`
	KeyCacheKeysToSave                            int64                            `yaml:"key_cache_keys_to_save,omitempty"`
	RowCacheClassName                             RowCacheClassName                `yaml:"row_cache_class_name,omitempty"`
	RowCacheSizeInMB                              int64                            `yaml:"row_cache_size_in_mb"`
	RowCacheSavePeriod                            int64                            `yaml:"row_cache_save_period"`
	RowCacheKeysToSave                            int64                            `yaml:"row_cache_keys_to_save,omitempty"`
	CounterCacheSizeInMB                          int64                            `yaml:"counter_cache_size_in_mb,omitempty"`
	CounterCacheSavePeriod                        int64                            `yaml:"counter_cache_save_period"`
	CounterCacheKeysToSave                        int64                            `yaml:"counter_cache_keys_to_save,omitempty"`
	SavedCachesDirectory                          string                           `yaml:"saved_caches_directory,omitempty"`
	CommitlogSync                                 commit.LogSync                   `yaml:"commitlog_sync"`
	CommitlogSyncBatchWindowInMS                  int64                            `yaml:"commitlog_sync_batch_window_in_ms,omitempty"`
	CommitlogSyncPeriodInMS                       int64                            `yaml:"commitlog_sync_period_in_ms"`
	CommitlogSegmentSizeInMB                      int64                            `yaml:"commitlog_segment_size_in_mb"`
	CommitlogCompression                          []Compression                    `yaml:"commitlog_compression,omitempty"`
	SeedProvider                                  []SeedProvider                   `yaml:"seed_provider"`
	ConcurrentReads                               int64                            `yaml:"concurrent_reads"`
	ConcurrentWrites                              int64                            `yaml:"concurrent_writes"`
	ConcurrentCounterWrites                       int64                            `yaml:"concurrent_counter_writes"`
	ConcurrentMaterializedViewWrites              int64                            `yaml:"concurrent_materialized_view_writes"`
	FileCacheSizeInMB                             int64                            `yaml:"file_cache_size_in_mb,omitempty"`
	BufferPoolUseHeapIfExhausted                  bool                             `yaml:"buffer_pool_use_heap_if_exhausted,omitempty"`
	DiskOptimizationStrategy                      disk.OptimizationStrategy        `yaml:"disk_optimization_strategy,omitempty"`
	MemtableHeapSpaceInMB                         int64                            `yaml:"memtable_heap_space_in_mb,omitempty"`
	MemtableOffheapSpaceInMB                      int64                            `yaml:"memtable_offheap_space_in_mb,omitempty"`
	MemtableCleanupThreshold                      float64                          `yaml:"memtable_cleanup_threshold,omitempty"`
	MemtableAllocationType                        MemtableAllocationType           `yaml:"memtable_allocation_type"`
	CommitlogTotalSpaceInMB                       int64                            `yaml:"commitlog_total_space_in_mb,omitempty"`
	MemtableFlushWriters                          int64                            `yaml:"memtable_flush_writers,omitempty"`
	CDCTotalSpaceInMB                             int64                            `yaml:"cdc_total_space_in_mb,omitempty"`
	CDCFreeSpaceCheckIntervalMS                   int64                            `yaml:"cdc_free_space_check_interval_ms,omitempty"`
	IndexSummaryCapacityInMB                      int64                            `yaml:"index_summary_capacity_in_mb,omitempty"`
	IndexSummaryResizeIntervalInMinutes           int64                            `yaml:"index_summary_resize_interval_in_minutes"`
	TrickleFSync                                  bool                             `yaml:"trickle_fsync"`
	TrickleFSyncIntervalInKB                      int64                            `yaml:"trickle_fsync_interval_in_kb"`
	StoragePort                                   int64                            `yaml:"storage_port"`
	StoragePortSSL                                int64                            `yaml:"ssl_storage_port"`
	ListenAddress                                 string                           `yaml:"listen_address"`
	ListenInterface                               string                           `yaml:"listen_interface,omitempty"`
	ListenInterfacePreferIPv6                     bool                             `yaml:"listen_interface_prefer_ipv6,omitempty"`
	BroadcastAddress                              string                           `yaml:"broadcast_address,omitempty"`
	ListenOnBroadcastAddress                      bool                             `yaml:"listen_on_broadcast_address,omitempty"`
	InternodeAuthenticator                        string                           `yaml:"internode_authenticator,omitempty"`
	StartNativeTransport                          bool                             `yaml:"start_native_transport"`
	NativeTransportPort                           int64                            `yaml:"native_transport_port"`
	NativeTransportPortSSL                        int64                            `yaml:"native_transport_port_ssl,omitempty"`
	NativeTransportMaxThreads                     int64                            `yaml:"native_transport_max_threads,omitempty"`
	NativeTransportMaxFrameSizeInMB               int64                            `yaml:"native_transport_max_frame_size_in_mb,omitempty"`
	NativeTransportMaxConcurrentConnections       int64                            `yaml:"native_transport_max_concurrent_connections,omitempty"`
	NativeTransportMaxConcurrentConnectionsPerIP  int64                            `yaml:"native_transport_max_concurrent_connections_per_ip,omitempty"`
	StartRPC                                      bool                             `yaml:"start_rpc"`
	RPCAddress                                    string                           `yaml:"rpc_address"`
	RPCInterface                                  string                           `yaml:"rpc_interface,omitempty"`
	RPCInterfacePreferIPv6                        bool                             `yaml:"rpc_interface_prefer_ipv6,omitempty"`
	RPCPort                                       int64                            `yaml:"rpc_port"`
	BroadcastRPCAddress                           string                           `yaml:"broadcast_rpc_address,omitempty"`
	RPCKeepalive                                  bool                             `yaml:"rpc_keepalive"`
	RPCServerType                                 RPCServerType                    `yaml:"rpc_server_type"`
	RPCMinThreads                                 int64                            `yaml:"rpc_min_threads,omitempty"`
	RPCMaxThreads                                 int64                            `yaml:"rpc_max_threads,omitempty"`
	RPCSendBuffSizeInBytes                        int64                            `yaml:"rpc_send_buff_size_in_bytes,omitempty"`
	RPCRecvBuffSizeInBytes                        int64                            `yaml:"rpc_recv_buff_size_in_bytes,omitempty"`
	InternodeSendBuffSizeInBytes                  int64                            `yaml:"internode_send_buff_size_in_bytes,omitempty"`
	InternodeRecvBuffSizeInBytes                  int64                            `yaml:"internode_recv_buff_size_in_bytes,omitempty"`
	ThriftFramedTransportSizeInMB                 int64                            `yaml:"thrift_framed_transport_size_in_mb"`
	IncrementalBackups                            bool                             `yaml:"incremental_backups"`
	SnapshotBeforeCompaction                      bool                             `yaml:"snapshot_before_compaction"`
	AutoSnapshot                                  bool                             `yaml:"auto_snapshot"`
	ColumnIndexSizeInKB                           int64                            `yaml:"column_index_size_in_kb"`
	ColumnIndexCacheSizeInKB                      int64                            `yaml:"column_index_cache_size_in_kb"`
	ConcurrentCompactors                          int64                            `yaml:"concurrent_compactors,omitempty"`
	CompactionThroughputMBPerSec                  int64                            `yaml:"compaction_throughput_mb_per_sec"`
	SSTablePreemptiveOpenIntervalInMB             int64                            `yaml:"sstable_preemptive_open_interval_in_mb"`
	StreamThroughputOutboundMegabitsPerSec        int64                            `yaml:"stream_throughput_outbound_megabits_per_sec,omitempty"`
	InterDCStreamThroughputOutboundMegabitsPerSec int64                            `yaml:"inter_dc_stream_throughput_outbound_megabits_per_sec,omitempty"`
	ReadRequestTimeoutInMS                        int64                            `yaml:"read_request_timeout_in_ms"`
	RangeRequestTimeoutInMS                       int64                            `yaml:"range_request_timeout_in_ms"`
	WriteRequestTimeoutInMS                       int64                            `yaml:"write_request_timeout_in_ms"`
	CounterWriteRequestTimeoutInMS                int64                            `yaml:"counter_write_request_timeout_in_ms"`
	CasContentionTimeoutInMS                      int64                            `yaml:"cas_contention_timeout_in_ms"`
	TruncateRequestTimeoutInMS                    int64                            `yaml:"truncate_request_timeout_in_ms"`
	RequestTimeoutInMS                            int64                            `yaml:"request_timeout_in_ms"`
	SlowQueryLogTimeoutInMS                       int64                            `yaml:"slow_query_log_timeout_in_ms"`
	CrossNodeTimeout                              bool                             `yaml:"cross_node_timeout"`
	StreamingKeepAlivePeriodInSecs                int64                            `yaml:"streaming_keep_alive_period_in_secs,omitempty"`
	StreamingConnectionsPerHost                   int64                            `yaml:"streaming_connections_per_host,omitempty"`
	PhiConvictThreshold                           float64                          `yaml:"phi_convict_threshold,omitempty"`
	EndpointSnitch                                string                           `yaml:"endpoint_snitch"`
	DynamicSnitchUpdateIntervalInMS               int64                            `yaml:"dynamic_snitch_update_interval_in_ms"`
	DynamicSnitchResetIntervalInMS                int64                            `yaml:"dynamic_snitch_reset_interval_in_ms"`
	DynamicSnitchBadnessThreshold                 float64                          `yaml:"dynamic_snitch_badness_threshold"`
	RequestScheduler                              string                           `yaml:"request_scheduler"`
	RequestSchedulerOptions                       RequestSchedulerOptions          `yaml:"request_scheduler_options,omitempty"`
	ServerEncryptionOptions                       ServerEncryptionOptions          `yaml:"server_encryption_options"`
	ClientEncryptionOptions                       ClientEncryptionOptions          `yaml:"client_encryption_options"`
	InternodeCompression                          InternodeCompression             `yaml:"internode_compression"`
	InterDCTCPNodelay                             bool                             `yaml:"inter_dc_tcp_nodelay"`
	TracetypeQueryTTL                             int64                            `yaml:"tracetype_query_ttl"`
	TracetypeRepairTTL                            int64                            `yaml:"tracetype_repair_ttl"`
	EnableUserDefinedFunctions                    bool                             `yaml:"enable_user_defined_functions"`
	EnableScriptedUserDefinedFunctions            bool                             `yaml:"enable_scripted_user_defined_functions"`
	EnableMaterializedViews                       bool                             `yaml:"enable_materialized_views"`
	WindowsTimerInterval                          int64                            `yaml:"windows_timer_interval"`
	TransparentDataEncryptionOptions              TransparentDataEncryptionOptions `yaml:"transparent_data_encryption_options"`
	TombstoneWarnThreshold                        int64                            `yaml:"tombstone_warn_threshold"`
	TombstoneFailureThreshold                     int64                            `yaml:"tombstone_failure_threshold"`
	BatchSizeWarnThresholdInKB                    int64                            `yaml:"batch_size_warn_threshold_in_kb"`
	BatchSizeFailThresholdInKB                    int64                            `yaml:"batch_size_fail_threshold_in_kb"`
	UnloggedBatchAcrossPartitionsWarnThreshold    int64                            `yaml:"unlogged_batch_across_partitions_warn_threshold"`
	CompactionLargePartitionWarningThresholdMB    int64                            `yaml:"compaction_large_partition_warning_threshold_mb"`
	GCLogThresholdInMS                            int64                            `yaml:"gc_log_threshold_in_ms,omitempty"`
	GCWarnThresholdInMS                           int64                            `yaml:"gc_warn_threshold_in_ms"`
	MaxValueSizeInMB                              int64                            `yaml:"max_value_size_in_mb,omitempty"`
	BackPressureEnabled                           bool                             `yaml:"back_pressure_enabled"`
	BackPressureStrategy                          []BackPressureStrategy           `yaml:"back_pressure_strategy"`
	OTCCoalescingStrategy                         string                           `yaml:"otc_coalescing_strategy,omitempty"`
	OTCCoalescingWindowUS                         int64                            `yaml:"otc_coalescing_window_us,omitempty"`
	OTCCoalescingEnoughCoalescedMessages          int64                            `yaml:"otc_coalescing_enough_coalesced_messages,omitempty"`
	OTCBacklogExpirationIntervalMS                int64                            `yaml:"otc_backlog_expiration_interval_ms,omitempty"`
	IdealConsistencyLevel                         string                           `yaml:"ideal_consistency_level,omitempty"`
}

// UnmarshalYAML sets default values and converts values from a YAML file.
func (c *CassandraYAML) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawCas CassandraYAML
	raw := rawCas{
		ClusterName:                "Test Cluster",
		NumTokens:                  256,
		HintedHandoffEnabled:       true,
		MaxHintWindowInMS:          10800000,
		HintedHandOffThrottleInKB:  1024,
		MaxHintsDeliveryThreads:    2,
		HintsFlushPeriodInMS:       10000,
		MaxHintsFileSizeInMB:       128,
		BatchlogReplayThrottleInKB: 1024,
		Authenticator:              "AllowAllAuthenticator",
		Authorizer:                 "AllowAllAuthorizer",
		RoleManager:                "CassandraRoleManager",
		RolesValidityInMS:          2000,
		PermissionsValidityInMS:    2000,
		CredentialsValidityInMS:    2000,
		Partitioner:                "org.apache.cassandra.dht.Murmur3Partitioner",
		CDCEnabled:                 false,
		DiskFailurePolicy:          disk.Stop,
		CommitFailurePolicy:        commit.Stop,
		KeyCacheSavePeriod:         14400,
		RowCacheSizeInMB:           0,
		RowCacheSavePeriod:         0,
		CounterCacheSavePeriod:     7200,
		CommitlogSync:              commit.Periodic,
		CommitlogSyncPeriodInMS:    10000,
		CommitlogSegmentSizeInMB:   32,
		SeedProvider: []SeedProvider{
			SeedProvider{
				ClassName: "org.apache.cassandra.locator.SimpleSeedProvider",
				Parameters: []SeedParameters{
					SeedParameters{Seeds: "127.0.0.1"},
				},
			},
		},
		ConcurrentReads:                     32,
		ConcurrentWrites:                    32,
		ConcurrentCounterWrites:             32,
		ConcurrentMaterializedViewWrites:    32,
		MemtableAllocationType:              HeapBuffers,
		IndexSummaryResizeIntervalInMinutes: 60,
		TrickleFSync:                        false,
		TrickleFSyncIntervalInKB:            10240,
		StoragePort:                         7000,
		StoragePortSSL:                      7001,
		ListenAddress:                       "localhost",
		StartNativeTransport:                true,
		NativeTransportPort:                 9042,
		StartRPC:                            false,
		RPCAddress:                          "localhost",
		RPCPort:                             9160,
		RPCKeepalive:                        true,
		RPCServerType:                       Sync,
		ThriftFramedTransportSizeInMB:       15,
		IncrementalBackups:                  false,
		SnapshotBeforeCompaction:            false,
		AutoSnapshot:                        true,
		ColumnIndexSizeInKB:                 64,
		ColumnIndexCacheSizeInKB:            2,
		CompactionThroughputMBPerSec:        16,
		SSTablePreemptiveOpenIntervalInMB:   50,
		ReadRequestTimeoutInMS:              5000,
		RangeRequestTimeoutInMS:             10000,
		WriteRequestTimeoutInMS:             2000,
		CounterWriteRequestTimeoutInMS:      5000,
		CasContentionTimeoutInMS:            1000,
		TruncateRequestTimeoutInMS:          60000,
		RequestTimeoutInMS:                  10000,
		SlowQueryLogTimeoutInMS:             500,
		CrossNodeTimeout:                    false,
		EndpointSnitch:                      "SimpleSnitch",
		DynamicSnitchUpdateIntervalInMS:     100,
		DynamicSnitchResetIntervalInMS:      600000,
		DynamicSnitchBadnessThreshold:       0.1,
		RequestScheduler:                    "org.apache.cassandra.scheduler.NoScheduler",
		ServerEncryptionOptions: ServerEncryptionOptions{
			InternodeEncryption: "none",
			Keystore:            "conf/.keystore",
			KeystorePassword:    "cassandra",
			Truststore:          "conf/.truststore",
			TruststorePassword:  "cassandra",
		},
		ClientEncryptionOptions: ClientEncryptionOptions{
			Enabled:          false,
			Optional:         false,
			Keystore:         "conf/.keystore",
			KeystorePassword: "cassandra",
		},
		InternodeCompression:               DC,
		InterDCTCPNodelay:                  false,
		TracetypeQueryTTL:                  86400,
		TracetypeRepairTTL:                 604800,
		EnableUserDefinedFunctions:         false,
		EnableScriptedUserDefinedFunctions: false,
		EnableMaterializedViews:            true,
		WindowsTimerInterval:               1,
		TransparentDataEncryptionOptions: TransparentDataEncryptionOptions{
			Enabled:       false,
			ChunkLengthKB: 64,
			Cipher:        "AES/CBC/PKCS5Padding",
			KeyAlias:      "testing:1",
			KeyProvider: []KeyProvider{
				KeyProvider{
					ClassName: "org.apache.cassandra.security.JKSKeyProvider",
					Parameters: []KeyProviderParameters{
						KeyProviderParameters{
							Keystore:         "conf/.keystore",
							KeystorePassword: "cassandra",
							StoreType:        "JCEKS",
							KeyPassword:      "cassandra",
						},
					},
				},
			},
		},
		TombstoneWarnThreshold:                     1000,
		TombstoneFailureThreshold:                  100000,
		BatchSizeWarnThresholdInKB:                 5,
		BatchSizeFailThresholdInKB:                 50,
		UnloggedBatchAcrossPartitionsWarnThreshold: 10,
		CompactionLargePartitionWarningThresholdMB: 100,
		GCWarnThresholdInMS:                        1000,
		BackPressureEnabled:                        false,
		BackPressureStrategy: []BackPressureStrategy{
			BackPressureStrategy{
				ClassName: "org.apache.cassandra.net.RateBasedBackPressure",
				Parameters: []map[string]interface{}{
					map[string]interface{}{
						"high_ratio": 0.90,
						"factor":     5,
						"flow":       "FAST",
					},
				},
			},
		},
	}
	if err := unmarshal(&raw); err != nil {
		return err
	}
	*c = CassandraYAML(raw)
	return nil
}
