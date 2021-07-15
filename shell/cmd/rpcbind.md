# rpcbind
主要功能是进行端口映射工作. 当客户端尝试连接并使用RPC服务器提供的服务（如NFS服务）时，rpcbind会将所管理的与服务对应的端口提供给客户端，从而使客户可以通过该端口向服务器请求服务.

    rpcbind中的每个program都有一个唯一与之对应的program number，它们的映射关系定义在/etc/rpc文件中.
    `rpcinfo -p localhost`:
    1. 第一列就是program number
    1. 第二列vers表示对应program的版本号, 即每个program可能启动了不同版本的功能
    1. 第三列是program监听的端口
    1. 最后一列为RPC管理的RPC service名, 其实就是各program对应的称呼

# rpcdebug
```bash
# rpcinfo -p 192.168.1.35 # 查看与nfs server间的rpc通信是否正常, service列表中显示mountd+nfs表示正常
# rpcdebug -vh
# rpcdebug -m nfs -s all # Enable all NFS (client-side) debugging
# rpcdebug -m rpc -s all # Enable RPC Call (client/server-side) debugging
# rpcdebug -m nfsd -s all # Enable NFSD (server-side) debugging
# ### Disable debugging
# rpcdebug -m nfs -c all
# rpcdebug -m nfsd -c all
```

rpcdebug module:
- nfs   NFS client
- nfsd  NFS server
- nlm   Network Lock Manager Protocol(NLM). 调试nfs锁管理器相关问题, 将只记录锁相关信息.
- rpc   Remote Procedure Call

rpcdebug选项:
- -m : module name to set or clear kernel debug flags
- -s : To set available kernel debug flag for a module
- -c : Clear Kernel debug flags
- -v : 显示更详细信息