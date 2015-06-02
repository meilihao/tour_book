#### `(error) MISCONF Redis is configured to save RDB snapshots, but is currently not able to persist on disk. Commands that may modify the data set are disabled. Please check Redis logs for details about the error.`

在redis-cli运行`ping`命令时碰到.
原因:强制停止redis快照导致不能持久化引起的(之前直接调用的github.com/garyburd/redigo/redis.Pool.Close()),具体原因可能是磁盘空间已满或持久化快照时没有权限(这种常见).

运行`config set stop-writes-on-bgsave-error no`命令或将`redis.conf # stop-writes-on-bgsave-error`设为no,可让redis忽略这个异常,不推荐.推荐是用适当权限重启redis-server.

#### 查看当前所有的key

运行`keys *`

#### `ERR Errors trying to SHUTDOWN. Check logs`

停止redis时碰到,可能是权限问题，因为在shutdown命令的时候，会进行save操作，而save需要操作dump.rdb文件，如果没有权限则会报这个错.

解决:先kill redis-server,再用适当权限(通常时sudo)重启.