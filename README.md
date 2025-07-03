# Restic Repo Exporter

```
RESTIC_PASSWORD="abc123" ./restic-repo-exporter --repo-path ./testdata
curl localhost:9100/metrics
# HELP go_gc_duration_seconds A summary of the wall-time pause (stop-the-world) duration in garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 0
go_gc_duration_seconds{quantile="0.25"} 0
go_gc_duration_seconds{quantile="0.5"} 0
go_gc_duration_seconds{quantile="0.75"} 0
go_gc_duration_seconds{quantile="1"} 0
go_gc_duration_seconds_sum 0
go_gc_duration_seconds_count 0
# HELP go_gc_gogc_percent Heap size target percentage configured by the user, otherwise 100. This value is set by the GOGC environment variable, and the runtime/debug.SetGCPercent function. Sourced from /gc/gogc:percent.
# TYPE go_gc_gogc_percent gauge
go_gc_gogc_percent 100
# HELP go_gc_gomemlimit_bytes Go runtime memory limit configured by the user, otherwise math.MaxInt64. This value is set by the GOMEMLIMIT environment variable, and the runtime/debug.SetMemoryLimit function. Sourced from /gc/gomemlimit:bytes.
# TYPE go_gc_gomemlimit_bytes gauge
go_gc_gomemlimit_bytes 9.223372036854776e+18
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 11
# HELP go_info Information about the Go environment.
# TYPE go_info gauge
go_info{version="go1.24.4"} 1
# HELP go_memstats_alloc_bytes Number of bytes allocated in heap and currently in use. Equals to /memory/classes/heap/objects:bytes.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 553120
# HELP go_memstats_alloc_bytes_total Total number of bytes allocated in heap until now, even if released already. Equals to /gc/heap/allocs:bytes.
# TYPE go_memstats_alloc_bytes_total counter
go_memstats_alloc_bytes_total 553120
# HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table. Equals to /memory/classes/profiling/buckets:bytes.
# TYPE go_memstats_buck_hash_sys_bytes gauge
go_memstats_buck_hash_sys_bytes 3701
# HELP go_memstats_frees_total Total number of heap objects frees. Equals to /gc/heap/frees:objects + /gc/heap/tiny/allocs:objects.
# TYPE go_memstats_frees_total counter
go_memstats_frees_total 0
# HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata. Equals to /memory/classes/metadata/other:bytes.
# TYPE go_memstats_gc_sys_bytes gauge
go_memstats_gc_sys_bytes 1.466048e+06
# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and currently in use, same as go_memstats_alloc_bytes. Equals to /memory/classes/heap/objects:bytes.
# TYPE go_memstats_heap_alloc_bytes gauge
go_memstats_heap_alloc_bytes 553120
# HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used. Equals to /memory/classes/heap/released:bytes + /memory/classes/heap/free:bytes.
# TYPE go_memstats_heap_idle_bytes gauge
go_memstats_heap_idle_bytes 466944
# HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use. Equals to /memory/classes/heap/objects:bytes + /memory/classes/heap/unused:bytes
# TYPE go_memstats_heap_inuse_bytes gauge
go_memstats_heap_inuse_bytes 3.072e+06
# HELP go_memstats_heap_objects Number of currently allocated objects. Equals to /gc/heap/objects:objects.
# TYPE go_memstats_heap_objects gauge
go_memstats_heap_objects 1178
# HELP go_memstats_heap_released_bytes Number of heap bytes released to OS. Equals to /memory/classes/heap/released:bytes.
# TYPE go_memstats_heap_released_bytes gauge
go_memstats_heap_released_bytes 466944
# HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system. Equals to /memory/classes/heap/objects:bytes + /memory/classes/heap/unused:bytes + /memory/classes/heap/released:bytes + /memory/classes/heap/free:bytes.
# TYPE go_memstats_heap_sys_bytes gauge
go_memstats_heap_sys_bytes 3.538944e+06
# HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
# TYPE go_memstats_last_gc_time_seconds gauge
go_memstats_last_gc_time_seconds 0
# HELP go_memstats_mallocs_total Total number of heap objects allocated, both live and gc-ed. Semantically a counter version for go_memstats_heap_objects gauge. Equals to /gc/heap/allocs:objects + /gc/heap/tiny/allocs:objects.
# TYPE go_memstats_mallocs_total counter
go_memstats_mallocs_total 1178
# HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures. Equals to /memory/classes/metadata/mcache/inuse:bytes.
# TYPE go_memstats_mcache_inuse_bytes gauge
go_memstats_mcache_inuse_bytes 14496
# HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system. Equals to /memory/classes/metadata/mcache/inuse:bytes + /memory/classes/metadata/mcache/free:bytes.
# TYPE go_memstats_mcache_sys_bytes gauge
go_memstats_mcache_sys_bytes 15704
# HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures. Equals to /memory/classes/metadata/mspan/inuse:bytes.
# TYPE go_memstats_mspan_inuse_bytes gauge
go_memstats_mspan_inuse_bytes 76320
# HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system. Equals to /memory/classes/metadata/mspan/inuse:bytes + /memory/classes/metadata/mspan/free:bytes.
# TYPE go_memstats_mspan_sys_bytes gauge
go_memstats_mspan_sys_bytes 81600
# HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place. Equals to /gc/heap/goal:bytes.
# TYPE go_memstats_next_gc_bytes gauge
go_memstats_next_gc_bytes 4.194304e+06
# HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations. Equals to /memory/classes/other:bytes.
# TYPE go_memstats_other_sys_bytes gauge
go_memstats_other_sys_bytes 1.143491e+06
# HELP go_memstats_stack_inuse_bytes Number of bytes obtained from system for stack allocator in non-CGO environments. Equals to /memory/classes/heap/stacks:bytes.
# TYPE go_memstats_stack_inuse_bytes gauge
go_memstats_stack_inuse_bytes 655360
# HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator. Equals to /memory/classes/heap/stacks:bytes + /memory/classes/os-stacks:bytes.
# TYPE go_memstats_stack_sys_bytes gauge
go_memstats_stack_sys_bytes 655360
# HELP go_memstats_sys_bytes Number of bytes obtained from system. Equals to /memory/classes/total:byte.
# TYPE go_memstats_sys_bytes gauge
go_memstats_sys_bytes 6.904848e+06
# HELP go_sched_gomaxprocs_threads The current runtime.GOMAXPROCS setting, or the number of operating system threads that can execute user-level Go code simultaneously. Sourced from /sched/gomaxprocs:threads.
# TYPE go_sched_gomaxprocs_threads gauge
go_sched_gomaxprocs_threads 12
# HELP go_threads Number of OS threads created.
# TYPE go_threads gauge
go_threads 11
# HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
# TYPE process_cpu_seconds_total counter
process_cpu_seconds_total 0
# HELP process_max_fds Maximum number of open file descriptors.
# TYPE process_max_fds gauge
process_max_fds 1.048576e+06
# HELP process_network_receive_bytes_total Number of bytes received by the process over the network.
# TYPE process_network_receive_bytes_total counter
process_network_receive_bytes_total 1.3522948338e+10
# HELP process_network_transmit_bytes_total Number of bytes sent by the process over the network.
# TYPE process_network_transmit_bytes_total counter
process_network_transmit_bytes_total 8.137965255e+09
# HELP process_open_fds Number of open file descriptors.
# TYPE process_open_fds gauge
process_open_fds 13
# HELP process_resident_memory_bytes Resident memory size in bytes.
# TYPE process_resident_memory_bytes gauge
process_resident_memory_bytes 1.1534336e+07
# HELP process_start_time_seconds Start time of the process since unix epoch in seconds.
# TYPE process_start_time_seconds gauge
process_start_time_seconds 1.75154167337e+09
# HELP process_virtual_memory_bytes Virtual memory size in bytes.
# TYPE process_virtual_memory_bytes gauge
process_virtual_memory_bytes 1.944809472e+09
# HELP process_virtual_memory_max_bytes Maximum amount of virtual memory available in bytes.
# TYPE process_virtual_memory_max_bytes gauge
process_virtual_memory_max_bytes 1.8446744073709552e+19
# HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
# TYPE promhttp_metric_handler_requests_in_flight gauge
promhttp_metric_handler_requests_in_flight 1
# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
# TYPE promhttp_metric_handler_requests_total counter
promhttp_metric_handler_requests_total{code="200"} 1
promhttp_metric_handler_requests_total{code="500"} 0
promhttp_metric_handler_requests_total{code="503"} 0
# HELP restic_repo_compression_progress_percent Compression progress of the repository in percent
# TYPE restic_repo_compression_progress_percent gauge
restic_repo_compression_progress_percent{repo="locked-repo"} 1
restic_repo_compression_progress_percent{repo="repo-with-tags"} 1
# HELP restic_repo_compression_ratio Compression ratio of the repository
# TYPE restic_repo_compression_ratio gauge
restic_repo_compression_ratio{repo="locked-repo"} 1.8274706867671693
restic_repo_compression_ratio{repo="repo-with-tags"} 1.393854748603352
# HELP restic_repo_compression_space_saving_percent Compression space saving of the repository in percent
# TYPE restic_repo_compression_space_saving_percent gauge
restic_repo_compression_space_saving_percent{repo="locked-repo"} 0.45279560036663613
restic_repo_compression_space_saving_percent{repo="repo-with-tags"} 0.282565130260521
# HELP restic_repo_last_snapshot_creation_seconds Time it took to create the last snapshot
# TYPE restic_repo_last_snapshot_creation_seconds gauge
restic_repo_last_snapshot_creation_seconds{hostname="Wortys-Thinkpad",repo="locked-repo",tag=""} 0.660879255
restic_repo_last_snapshot_creation_seconds{hostname="test-host",repo="repo-with-tags",tag="test-tag"} 0.723240066
restic_repo_last_snapshot_creation_seconds{hostname="test-host",repo="repo-with-tags",tag="test-tag-latest"} 0.640371108
# HELP restic_repo_last_snapshot_timestamp Timestamp of the last snapshot in the repository by hostname and tag
# TYPE restic_repo_last_snapshot_timestamp gauge
restic_repo_last_snapshot_timestamp{hostname="Wortys-Thinkpad",repo="locked-repo",tag=""} 1.75149266e+09
restic_repo_last_snapshot_timestamp{hostname="test-host",repo="repo-with-tags",tag="test-tag"} 1.751478405e+09
restic_repo_last_snapshot_timestamp{hostname="test-host",repo="repo-with-tags",tag="test-tag-latest"} 1.751478413e+09
# HELP restic_repo_num_errors Total number of errors found in the repository during check
# TYPE restic_repo_num_errors gauge
restic_repo_num_errors{repo="locked-repo"} 0
restic_repo_num_errors{repo="repo-with-tags"} 0
# HELP restic_repo_number_of_snapshots Total number of snapshots in the repository by hostname and tag
# TYPE restic_repo_number_of_snapshots gauge
restic_repo_number_of_snapshots{hostname="Wortys-Thinkpad",repo="locked-repo",tag=""} 2
restic_repo_number_of_snapshots{hostname="test-host",repo="repo-with-tags",tag="test-tag"} 1
restic_repo_number_of_snapshots{hostname="test-host",repo="repo-with-tags",tag="test-tag-latest"} 1
# HELP restic_repo_scrape_duration_seconds 
# TYPE restic_repo_scrape_duration_seconds histogram
restic_repo_scrape_duration_seconds_bucket{action="check",repo="locked-repo",le="0.1"} 0
restic_repo_scrape_duration_seconds_bucket{action="check",repo="locked-repo",le="0.20355579570665744"} 0
restic_repo_scrape_duration_seconds_bucket{action="check",repo="locked-repo",le="0.41434961965770456"} 1
restic_repo_scrape_duration_seconds_bucket{action="check",repo="locked-repo",le="0.8434326653017492"} 1
restic_repo_scrape_duration_seconds_bucket{action="check",repo="locked-repo",le="1.7168560731048441"} 1
restic_repo_scrape_duration_seconds_bucket{action="check",repo="locked-repo",le="3.4947600407466375"} 1
restic_repo_scrape_duration_seconds_bucket{action="check",repo="locked-repo",le="7.1137866089801225"} 1
restic_repo_scrape_duration_seconds_bucket{action="check",repo="locked-repo",le="14.480524936783134"} 1
restic_repo_scrape_duration_seconds_bucket{action="check",repo="locked-repo",le="29.475947757569855"} 1
restic_repo_scrape_duration_seconds_bucket{action="check",repo="locked-repo",le="59.99999999999997"} 1
restic_repo_scrape_duration_seconds_bucket{action="check",repo="locked-repo",le="+Inf"} 1
restic_repo_scrape_duration_seconds_sum{action="check",repo="locked-repo"} 0.338797894
restic_repo_scrape_duration_seconds_count{action="check",repo="locked-repo"} 1
restic_repo_scrape_duration_seconds_bucket{action="check",repo="repo-with-tags",le="0.1"} 0
restic_repo_scrape_duration_seconds_bucket{action="check",repo="repo-with-tags",le="0.20355579570665744"} 0
restic_repo_scrape_duration_seconds_bucket{action="check",repo="repo-with-tags",le="0.41434961965770456"} 1
restic_repo_scrape_duration_seconds_bucket{action="check",repo="repo-with-tags",le="0.8434326653017492"} 1
restic_repo_scrape_duration_seconds_bucket{action="check",repo="repo-with-tags",le="1.7168560731048441"} 1
restic_repo_scrape_duration_seconds_bucket{action="check",repo="repo-with-tags",le="3.4947600407466375"} 1
restic_repo_scrape_duration_seconds_bucket{action="check",repo="repo-with-tags",le="7.1137866089801225"} 1
restic_repo_scrape_duration_seconds_bucket{action="check",repo="repo-with-tags",le="14.480524936783134"} 1
restic_repo_scrape_duration_seconds_bucket{action="check",repo="repo-with-tags",le="29.475947757569855"} 1
restic_repo_scrape_duration_seconds_bucket{action="check",repo="repo-with-tags",le="59.99999999999997"} 1
restic_repo_scrape_duration_seconds_bucket{action="check",repo="repo-with-tags",le="+Inf"} 1
restic_repo_scrape_duration_seconds_sum{action="check",repo="repo-with-tags"} 0.340362295
restic_repo_scrape_duration_seconds_count{action="check",repo="repo-with-tags"} 1
restic_repo_scrape_duration_seconds_bucket{action="raw-stats",repo="locked-repo",le="0.1"} 0
restic_repo_scrape_duration_seconds_bucket{action="raw-stats",repo="locked-repo",le="0.20355579570665744"} 0
restic_repo_scrape_duration_seconds_bucket{action="raw-stats",repo="locked-repo",le="0.41434961965770456"} 1
restic_repo_scrape_duration_seconds_bucket{action="raw-stats",repo="locked-repo",le="0.8434326653017492"} 1
restic_repo_scrape_duration_seconds_bucket{action="raw-stats",repo="locked-repo",le="1.7168560731048441"} 1
restic_repo_scrape_duration_seconds_bucket{action="raw-stats",repo="locked-repo",le="3.4947600407466375"} 1
restic_repo_scrape_duration_seconds_bucket{action="raw-stats",repo="locked-repo",le="7.1137866089801225"} 1
restic_repo_scrape_duration_seconds_bucket{action="raw-stats",repo="locked-repo",le="14.480524936783134"} 1
restic_repo_scrape_duration_seconds_bucket{action="raw-stats",repo="locked-repo",le="29.475947757569855"} 1
restic_repo_scrape_duration_seconds_bucket{action="raw-stats",repo="locked-repo",le="59.99999999999997"} 1
restic_repo_scrape_duration_seconds_bucket{action="raw-stats",repo="locked-repo",le="+Inf"} 1
restic_repo_scrape_duration_seconds_sum{action="raw-stats",repo="locked-repo"} 0.287787078
restic_repo_scrape_duration_seconds_count{action="raw-stats",repo="locked-repo"} 1
restic_repo_scrape_duration_seconds_bucket{action="raw-stats",repo="repo-with-tags",le="0.1"} 0
restic_repo_scrape_duration_seconds_bucket{action="raw-stats",repo="repo-with-tags",le="0.20355579570665744"} 0
restic_repo_scrape_duration_seconds_bucket{action="raw-stats",repo="repo-with-tags",le="0.41434961965770456"} 1
restic_repo_scrape_duration_seconds_bucket{action="raw-stats",repo="repo-with-tags",le="0.8434326653017492"} 1
restic_repo_scrape_duration_seconds_bucket{action="raw-stats",repo="repo-with-tags",le="1.7168560731048441"} 1
restic_repo_scrape_duration_seconds_bucket{action="raw-stats",repo="repo-with-tags",le="3.4947600407466375"} 1
restic_repo_scrape_duration_seconds_bucket{action="raw-stats",repo="repo-with-tags",le="7.1137866089801225"} 1
restic_repo_scrape_duration_seconds_bucket{action="raw-stats",repo="repo-with-tags",le="14.480524936783134"} 1
restic_repo_scrape_duration_seconds_bucket{action="raw-stats",repo="repo-with-tags",le="29.475947757569855"} 1
restic_repo_scrape_duration_seconds_bucket{action="raw-stats",repo="repo-with-tags",le="59.99999999999997"} 1
restic_repo_scrape_duration_seconds_bucket{action="raw-stats",repo="repo-with-tags",le="+Inf"} 1
restic_repo_scrape_duration_seconds_sum{action="raw-stats",repo="repo-with-tags"} 0.297011889
restic_repo_scrape_duration_seconds_count{action="raw-stats",repo="repo-with-tags"} 1
restic_repo_scrape_duration_seconds_bucket{action="snapshots",repo="locked-repo",le="0.1"} 0
restic_repo_scrape_duration_seconds_bucket{action="snapshots",repo="locked-repo",le="0.20355579570665744"} 0
restic_repo_scrape_duration_seconds_bucket{action="snapshots",repo="locked-repo",le="0.41434961965770456"} 1
restic_repo_scrape_duration_seconds_bucket{action="snapshots",repo="locked-repo",le="0.8434326653017492"} 1
restic_repo_scrape_duration_seconds_bucket{action="snapshots",repo="locked-repo",le="1.7168560731048441"} 1
restic_repo_scrape_duration_seconds_bucket{action="snapshots",repo="locked-repo",le="3.4947600407466375"} 1
restic_repo_scrape_duration_seconds_bucket{action="snapshots",repo="locked-repo",le="7.1137866089801225"} 1
restic_repo_scrape_duration_seconds_bucket{action="snapshots",repo="locked-repo",le="14.480524936783134"} 1
restic_repo_scrape_duration_seconds_bucket{action="snapshots",repo="locked-repo",le="29.475947757569855"} 1
restic_repo_scrape_duration_seconds_bucket{action="snapshots",repo="locked-repo",le="59.99999999999997"} 1
restic_repo_scrape_duration_seconds_bucket{action="snapshots",repo="locked-repo",le="+Inf"} 1
restic_repo_scrape_duration_seconds_sum{action="snapshots",repo="locked-repo"} 0.283024765
restic_repo_scrape_duration_seconds_count{action="snapshots",repo="locked-repo"} 1
restic_repo_scrape_duration_seconds_bucket{action="snapshots",repo="repo-with-tags",le="0.1"} 0
restic_repo_scrape_duration_seconds_bucket{action="snapshots",repo="repo-with-tags",le="0.20355579570665744"} 0
restic_repo_scrape_duration_seconds_bucket{action="snapshots",repo="repo-with-tags",le="0.41434961965770456"} 1
restic_repo_scrape_duration_seconds_bucket{action="snapshots",repo="repo-with-tags",le="0.8434326653017492"} 1
restic_repo_scrape_duration_seconds_bucket{action="snapshots",repo="repo-with-tags",le="1.7168560731048441"} 1
restic_repo_scrape_duration_seconds_bucket{action="snapshots",repo="repo-with-tags",le="3.4947600407466375"} 1
restic_repo_scrape_duration_seconds_bucket{action="snapshots",repo="repo-with-tags",le="7.1137866089801225"} 1
restic_repo_scrape_duration_seconds_bucket{action="snapshots",repo="repo-with-tags",le="14.480524936783134"} 1
restic_repo_scrape_duration_seconds_bucket{action="snapshots",repo="repo-with-tags",le="29.475947757569855"} 1
restic_repo_scrape_duration_seconds_bucket{action="snapshots",repo="repo-with-tags",le="59.99999999999997"} 1
restic_repo_scrape_duration_seconds_bucket{action="snapshots",repo="repo-with-tags",le="+Inf"} 1
restic_repo_scrape_duration_seconds_sum{action="snapshots",repo="repo-with-tags"} 0.277526992
restic_repo_scrape_duration_seconds_count{action="snapshots",repo="repo-with-tags"} 1
# HELP restic_repo_suggest_prune Whether the repository suggests pruning
# TYPE restic_repo_suggest_prune gauge
restic_repo_suggest_prune{repo="locked-repo"} 0
restic_repo_suggest_prune{repo="repo-with-tags"} 0
# HELP restic_repo_suggest_repair_index Whether the repository suggests repairing the index
# TYPE restic_repo_suggest_repair_index gauge
restic_repo_suggest_repair_index{repo="locked-repo"} 0
restic_repo_suggest_repair_index{repo="repo-with-tags"} 0
# HELP restic_repo_total_blob_count Total number of blobs in the repository
# TYPE restic_repo_total_blob_count gauge
restic_repo_total_blob_count{repo="locked-repo"} 2
restic_repo_total_blob_count{repo="repo-with-tags"} 2
# HELP restic_repo_total_size_bytes Total size of the repository in bytes
# TYPE restic_repo_total_size_bytes gauge
restic_repo_total_size_bytes{repo="locked-repo"} 597
restic_repo_total_size_bytes{repo="repo-with-tags"} 358
# HELP restic_repo_total_snapshots_count Total number of snapshots in the repository
# TYPE restic_repo_total_snapshots_count gauge
restic_repo_total_snapshots_count{repo="locked-repo"} 2
restic_repo_total_snapshots_count{repo="repo-with-tags"} 2
# HELP restic_repo_total_uncompressed_size_bytes Total uncompressed size of the repository in bytes
# TYPE restic_repo_total_uncompressed_size_bytes gauge
restic_repo_total_uncompressed_size_bytes{repo="locked-repo"} 1091
restic_repo_total_uncompressed_size_bytes{repo="repo-with-tags"} 499
```
