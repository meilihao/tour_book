## gogs

### FAQ

#### `ssh: connect to host $host_name port 22: Connection refused`

git clone时碰到,原因是ssh服务未运行.

    sudo apt-get install openssh-server
    sudo systemctl start ssh