box.cfg{
    listen              = 3304,
    username            = "tarantool",
    memtx_memory        = 1*256*1024*1024,
    slab_alloc_factor   = 1.04,
    memtx_max_tuple_size= 8388608,
    --rows_per_wal        = 100000,
    custom_proc_title   = "heckfy_shard1",
    work_dir            = "/var/tarantool_heckfy_shard1",
    wal_dir             = "/var/tarantool_heckfy_shard1/xlogs",
    memtx_dir           = "/var/tarantool_heckfy_shard1/snaps",
    read_only           = false,
    pid_file            = "box.pid",
    log                 = "/var/tarantool_heckfy_shard1/logs/tarantool_heckfy_shard1.log",
    background          = true,     
    checkpoint_interval = 0,
    readahead = 1048576,
}
require('console').listen('127.0.0.1:3314')

dofile('/var/tarantool_heckfy_shard1/heckfy.lua')
dofile('/var/tarantool_heckfy_shard1/tarantool_heckfy.grants.lua')


box.once("schema", function()
    --makemegrants()
    print('box.once executed on master')
end)
